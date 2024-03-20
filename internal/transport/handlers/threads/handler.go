package threads

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/handlers/common"
	commonviews "goBoard/internal/transport/handlers/common/views"
	"goBoard/internal/transport/handlers/threads/views"
	"goBoard/internal/transport/middlewares/jwtauth"
	"goBoard/internal/transport/middlewares/session"
	"strings"
	"time"

	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	ElitismTitle = "elitism. secrecy. tradition."
)

type Handler struct {
	threadService ports.ThreadService
	memberService ports.MemberService

	jwtAuth *jwtauth.JWTAuth

	logger *zap.SugaredLogger
}

func NewHandler(threadService ports.ThreadService, memberService ports.MemberService, jwtAuth *jwtauth.JWTAuth, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		threadService: threadService,
		memberService: memberService,
		jwtAuth:       jwtAuth,
		logger:        logger,
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(h.jwtAuth))
		r.Use(jwtauth.Authenticator(h.jwtAuth))

		r.Get("/", h.ThreadsHome)
		r.Get("/threads", h.Threads)
		r.Get("/thread/view/{id}", h.Thread)
		r.Post("/thread/create", h.CreateThread)
		r.Get("/thread/create", h.NewThreadPage)
		r.Post("/preview", h.Preview)
		r.Post("/post", h.Post)
		r.Get("/posts", h.Posts)
		r.Get("/dot/{threadId}", h.ToggleDot)
	})
}

func (h *Handler) NewThreadPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	templ.Handler(views.Home(views.NewThreadForm(cookie.Value), views.NewThreadTitleGroup(), cookie.Value)).Component.Render(r.Context(), w)
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	sess, err := session.Get("member", r)
	if err != nil {
		h.logger.Errorf("error getting session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username, ok := sess.Values["name"]
	if !ok {
		h.logger.Errorf("username is not in session")
		w.WriteHeader(http.StatusInternalServerError)
	}

	usernameStr, ok := username.(string)
	if !ok {
		h.logger.Errorf("username is not a string: %v", username)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subject := r.FormValue("subject")
	if subject == "" {
		h.logger.Error("subject is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := r.FormValue("body")
	if body == "" {
		h.logger.Error("body is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ip := r.RemoteAddr
	if strings.Contains(ip, "[::1]") {
		ip = "127.0.0.1"
	}

	threadID, err := h.threadService.NewThread(usernameStr, ip, body, subject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/thread/view/"+strconv.Itoa(threadID), http.StatusFound)
}

func (h *Handler) ToggleDot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sess, err := session.Get("member", r)
	if err != nil {
		h.logger.Errorf("error getting session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memberID, ok := sess.Values["id"]
	if !ok {
		h.logger.Errorf("memberID is not in session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memberIDInt, ok := memberID.(int)
	if !ok {
		h.logger.Errorf("memberID is not an int: %v", memberID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	threadID := chi.URLParam(r, "threadId")
	if threadID == "" {
		h.logger.Error("threadID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	threadIDInt, err := strconv.Atoi(threadID)
	if err != nil {
		h.logger.Errorf("threadID is not an int: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dotted, err := h.threadService.ToggleDot(ctx, memberIDInt, threadIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(commonviews.DotControl(dotted, threadIDInt)).Component.Render(ctx, w)
}

func (h *Handler) Posts(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.logger.Errorf("error getting cookie: %v", err)
	}

	sess, err := session.Get("member", r)
	if err != nil {
		h.logger.Errorf("error getting session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	viewable, ok := sess.Values["collapseopen"]
	if !ok {
		viewable = 20
	}

	viewableInt, ok := viewable.(int)
	if !ok {
		h.logger.Errorf("viewable is not an int: %v", viewable)
		w.WriteHeader(http.StatusInternalServerError)
	}

	memberID, ok := sess.Values["id"]
	if !ok {
		h.logger.Errorf("memberID is not in session")
		w.WriteHeader(http.StatusInternalServerError)
	}

	memberIDInt, ok := memberID.(int)
	if !ok {
		h.logger.Errorf("memberID is not an int: %v", memberID)
		w.WriteHeader(http.StatusInternalServerError)
	}

	start, err := strconv.Atoi(r.URL.Query().Get("start"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	threadID, err := strconv.Atoi(r.URL.Query().Get("threadId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	end, err := strconv.Atoi(r.URL.Query().Get("end"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	thread, err := h.threadService.GetCollapsibleThreadByID(r.Context(), end-start, threadID, memberIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(commonviews.Posts(common.ThreadToPosts(*thread), viewableInt, cookie.Value)).Component.Render(r.Context(), w)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	sess, err := session.Get("member", r)
	if err != nil {
		h.logger.Errorf("error getting session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username, ok := sess.Values["name"]
	if !ok {
		h.logger.Errorf("username is not in session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		h.logger.Errorf("username is not a string: %v", username)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := r.FormValue("body")
	if body == "" {
		h.logger.Error("body is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	threadID := r.FormValue("threadId")
	if threadID == "" {
		h.logger.Error("threadID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idx := r.FormValue("idx")
	if idx == "" {
		h.logger.Error("idx is empty")
		w.WriteHeader(http.StatusBadRequest)
	}

	idxInt, err := strconv.Atoi(idx)
	if err != nil {
		h.logger.Errorf("idx is not an int: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	threadIDInt, err := strconv.Atoi(threadID)
	if err != nil {
		h.logger.Errorf("threadID is not an int: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ip := r.RemoteAddr
	if strings.Contains(ip, "[::1]") {
		ip = "127.0.0.1"
	}

	_, err = h.threadService.NewPost(body, ip, usernameStr, threadIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post := common.Post{
		Body:       body,
		MemberName: usernameStr,
		Date:       time.Now().Format("Mon Jan 2, 2006 03:04 pm"),
	}

	templ.Handler(commonviews.Post(post, idxInt, true)).Component.Render(r.Context(), w)
}

func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil {
		h.logger.Errorf("error getting cookie: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
	}

	err = r.ParseForm()
	if err != nil {
		h.logger.Errorf("error parsing form: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := r.FormValue("body")
	if body == "" {
		h.logger.Error("body is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idx := r.FormValue("idx")
	if idx == "" {
		h.logger.Error("idx is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idxInt, err := strconv.Atoi(idx)
	if err != nil {
		h.logger.Errorf("idx is not an int: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Format("Mon Jan 2, 2006 03:04 pm")

	post := common.Post{
		Body:       body,
		MemberName: cookie.Value,
		Date:       timestamp,
		Preview:    true,
	}

	templ.Handler(commonviews.Post(post, idxInt, true)).Component.Render(r.Context(), w)
}

func (h *Handler) Thread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("user")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sess, err := session.Get("member", r)
	if err != nil {
		h.logger.Errorf("error getting session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	viewable, ok := sess.Values["collapseopen"]
	if !ok {
		viewable = 20
	}

	viewableInt, ok := viewable.(int)
	if !ok {
		h.logger.Errorf("viewable is not an int: %v", viewable)
		w.WriteHeader(http.StatusInternalServerError)
	}

	memberID, ok := sess.Values["id"]
	if !ok {
		h.logger.Errorf("memberID is not in session")
		w.WriteHeader(http.StatusInternalServerError)
	}

	memberIDInt, ok := memberID.(int)
	if !ok {
		h.logger.Errorf("memberID is not an int: %v", memberID)
		w.WriteHeader(http.StatusInternalServerError)
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	thread, err := h.threadService.GetCollapsibleThreadByID(ctx, viewableInt, idInt, memberIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	posts := common.ThreadToPosts(*thread)

	templ.Handler(views.Home(views.Thread(commonviews.PostsPage(posts, viewableInt, cookie.Value), thread.ID), commonviews.PostsTitleGroup(*thread), cookie.Value)).Component.Render(ctx, w)
}

func (h *Handler) ThreadsHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("user")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := session.Get("member", r)
	if err != nil {
		h.logger.Errorf("error getting session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memberID, ok := session.Values["id"]
	if !ok {
		h.logger.Errorf("memberID is not in session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memberIDInt, ok := memberID.(int)
	if !ok {
		h.logger.Errorf("memberID is not an int: %v", memberID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	threads, cursorsOut, err := h.threadService.ListThreads(ctx, domain.Cursors{}, 50, memberIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(views.Home(views.Threads(threads, cursorsOut), views.ThreadsTitleGroup(ElitismTitle), cookie.Value)).Component.Render(ctx, w)
}

func (h *Handler) Threads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("user")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := session.Get("member", r)

	memberID, ok := session.Values["id"]
	if !ok {
		h.logger.Errorf("memberID is not in session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memberIDInt, ok := memberID.(int)
	if !ok {
		h.logger.Errorf("memberID is not an int: %v", memberID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	nextCursor := r.URL.Query().Get("next")
	prevCursor := r.URL.Query().Get("prev")
	dir := r.URL.Query().Get("dir")
	if dir != "next" && dir != "prev" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if dir == "next" {
		prevCursor = ""
	} else {
		nextCursor = ""
	}

	cursors := domain.Cursors{
		Next: nextCursor,
		Prev: prevCursor,
	}

	threads, cursors, err := h.threadService.ListThreads(ctx, cursors, 50, memberIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isHx(r) {
		templ.Handler(views.Threads(threads, cursors)).Component.Render(ctx, w)
		return
	}

	templ.Handler((views.Home(views.Threads(threads, cursors), views.ThreadsTitleGroup(ElitismTitle), cookie.Value))).Component.Render(ctx, w)
}

func isHx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
