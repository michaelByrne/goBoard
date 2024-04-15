package authentication

import (
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/handlers/authentication/views"
	"goBoard/internal/transport/middlewares/session"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	authService ports.AuthenticationService
}

func NewHandler(authService ports.AuthenticationService) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/login", h.Login)
	r.Get("/login", h.LoginForm)
	r.Get("/logout/{name}", h.Logout)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	setTokenCookie("access-token", "", time.Now(), w)

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	templ.Handler(views.Login()).Component.Render(r.Context(), w)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := r.FormValue("name")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password := r.FormValue("pass")
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	member, token, err := h.authService.Authenticate(ctx, username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if member == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	setTokenCookie("access-token", token.TokenStr, token.Expires, w)

	sess, err := session.Get("member", r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sess.Values["id"] = member.ID
	sess.Values["name"] = username
	sess.Values["admin"] = member.IsAdmin

	for key, pref := range member.Prefs {
		sess.Values[key] = pref.Value
	}

	err = sess.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func setTokenCookie(name, token string, expiration time.Time, w http.ResponseWriter) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	http.SetCookie(w, cookie)
}
