package members

import (
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/handlers/common"
	"goBoard/internal/transport/handlers/members/views"
	threadviews "goBoard/internal/transport/handlers/threads/views"
	"goBoard/internal/transport/middlewares/jwtauth"
	"goBoard/internal/transport/middlewares/session"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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

		r.Get("/profile/{username}", h.Profile)
	})
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	username := chi.URLParam(r, "username")
	if username == "" {
		sess, err := session.Get("member", r)
		if err != nil {
			h.logger.Errorf("error getting session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		member, err := common.GetMember(*sess)
		if err != nil {
			h.logger.Errorf("error getting member: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		username = member.Username
	}

	member, err := h.memberService.GetMemberByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(threadviews.Home(views.Profile(*member), views.ProfileTitleGroup(*member), member.Name)).Component.Render(ctx, w)
}
