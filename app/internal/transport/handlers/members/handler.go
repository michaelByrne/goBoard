package members

import (
	"errors"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/gberrors"
	"goBoard/internal/transport/handlers/common"
	"goBoard/internal/transport/handlers/members/views"
	messageviews "goBoard/internal/transport/handlers/messages/views"
	threadviews "goBoard/internal/transport/handlers/threads/views"
	"goBoard/internal/transport/middlewares/session"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	threadService ports.ThreadService
	memberService ports.MemberService

	verifyMiddleware func(next http.Handler) http.Handler

	logger *zap.SugaredLogger
}

func NewHandler(threadService ports.ThreadService, memberService ports.MemberService, verifyMiddleware func(next http.Handler) http.Handler, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		threadService:    threadService,
		memberService:    memberService,
		verifyMiddleware: verifyMiddleware,
		logger:           logger,
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(h.verifyMiddleware)

		r.Get("/profile/{username}", h.Profile)
		r.Get("/validate", h.Validate)
		r.Get("/prefs", h.Prefs)
		r.Get("/member/edit", h.EditPage)
		r.Post("/member/edit", h.Edit)
	})
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sess, err := session.Get("member", r)
	if err != nil {
		h.logger.Errorf("error getting session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionMember, err := common.GetMember(*sess)
	if err != nil {
		h.logger.Errorf("error getting member: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.logger.Errorf("error parsing form: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newPrefs := make(map[string]domain.MemberPref)
	var postal string
	for k, v := range r.Form {
		if len(v) == 0 || v[0] == "" || k == "username" {
			continue
		}

		if len(v) > 1 {
			h.logger.Errorf("multiple values for key: %s", k)
			continue
		}

		if k == "postal" {
			err = h.memberService.UpdatePostalCode(ctx, sessionMember.ID, v[0])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			postal = v[0]

			continue
		}

		newPrefs[k] = domain.MemberPref{
			Value: v[0],
		}
	}

	err = h.memberService.UpdatePrefs(ctx, sessionMember.ID, newPrefs)
	if err != nil {
		h.logger.Errorf("error updating prefs: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prefs, err := h.memberService.GetMergedPrefs(ctx, sessionMember.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(views.PrefsWithSwap(prefs, postal)).Component.Render(ctx, w)
}

func (h *Handler) EditPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sess, err := session.Get("member", r)
	if err != nil {
		h.logger.Errorf("error getting session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionMember, err := common.GetMember(*sess)
	if err != nil {
		h.logger.Errorf("error getting member: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	member, err := h.memberService.GetMemberByID(sessionMember.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prefs, err := h.memberService.GetMergedPrefs(ctx, member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(threadviews.Home(views.EditProfile(*member, prefs), views.AccountManagementTitleGroup(member.Name), member.Name)).Component.Render(ctx, w)
}

func (h *Handler) Prefs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	prefs, err := h.memberService.GetMergedPrefs(ctx, member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(views.Prefs(prefs)).Component.Render(ctx, w)
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sess, err := session.Get("member", r)
	if err != nil {
		h.logger.Errorf("error getting session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionMember, err := common.GetMember(*sess)
	if err != nil {
		h.logger.Errorf("error getting member: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	namesStr := r.URL.Query().Get("names")
	if namesStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	names := strings.Split(namesStr, ",")

	var members []domain.Member
	for _, name := range names {
		member, err := h.memberService.GetMemberByUsername(name)
		if err != nil {
			var notFoundErr gberrors.MemberNotFound
			if errors.As(err, &notFoundErr) {
				continue
			}

			h.logger.Errorf("error getting member by username: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if sessionMember.Username != name {
			members = append(members, *member)
		}
	}

	templ.Handler(messageviews.Recipients(members)).Component.Render(ctx, w)
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
