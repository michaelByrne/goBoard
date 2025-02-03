package threads

import (
	"errors"
	"fmt"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/gberrors"
	"goBoard/internal/transport/handlers/common"
	commonviews "goBoard/internal/transport/handlers/common/views"
	"goBoard/internal/transport/handlers/threads/views"
	"goBoard/internal/transport/middlewares/session"
	"io"
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
	imageService  ports.ImageService
	themeService  ports.ThemeService

	verifyMiddleware func(next http.Handler) http.Handler

	logger *zap.SugaredLogger
}

func NewHandler(threadService ports.ThreadService, memberService ports.MemberService, imageService ports.ImageService, themeService ports.ThemeService, verifyMiddleware func(next http.Handler) http.Handler, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		threadService:    threadService,
		memberService:    memberService,
		imageService:     imageService,
		themeService:     themeService,
		verifyMiddleware: verifyMiddleware,
		logger:           logger,
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(h.verifyMiddleware)

		r.Get("/", h.ThreadsHome)
		r.Get("/threads", h.Threads)
		r.Get("/ignored/{memberId}", h.IgnoredThreads)
		r.Get("/created/{memberId}", h.CreatedThreads)
		r.Get("/participated/{memberId}", h.ParticipatedThreads)
		r.Get("/favorited/{memberId}", h.FavoritedThreads)
		r.Get("/thread/view/{id}", h.Thread)
		r.Get("/view/thread/{threadId}", h.ViewThread)
		r.Post("/thread/create", h.CreateThread)
		r.Get("/thread/create", h.NewThreadPage)
		r.Post("/thread/preview", h.Preview)
		r.Post("/thread/post", h.Post)
		r.Get("/thread/posts", h.Posts)
		r.Get("/dot/{threadId}", h.ToggleDot)
		r.Get("/ignore/{threadId}", h.ToggleIgnore)
		r.Get("/favorite/{threadId}", h.ToggleFavorite)
		r.Get("/uploader", h.Uploader)
		r.Post("/image/upload", h.UploadImage)
		r.Get("/image/refresh/{key}", h.RefreshImage)
		r.Get("/styles", h.Styles)
	})
}

func (h *Handler) Styles(w http.ResponseWriter, r *http.Request) {
	//sess, err := session.Get("member", r)
	//if err != nil {
	//	h.logger.Errorf("error getting session: %v", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//member, err := common.GetMember(*sess)
	//if err != nil {
	//	h.logger.Errorf("error getting member: %v", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//}

	theme, err := h.themeService.GetTheme(r.Context(), "blue")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(views.Styles(theme)).Component.Render(r.Context(), w)
}

func (h *Handler) RefreshImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key := chi.URLParam(r, "key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	presignedURL, err := h.imageService.RefreshPresign(ctx, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(commonviews.InitialImage(presignedURL, key)).Component.Render(ctx, w)
}

func (h *Handler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
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

	isFavorite, err := h.threadService.ToggleFavorite(ctx, member.ID, threadIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(commonviews.FavoriteControl(isFavorite, threadIDInt)).Component.Render(ctx, w)
}

func (h *Handler) ViewThread(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.threadService.ViewThread(ctx, member.ID, threadIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ToggleIgnore(w http.ResponseWriter, r *http.Request) {
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

	ignored, err := h.threadService.ToggleIgnore(ctx, member.ID, threadIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(commonviews.IgnoreControl(ignored, threadIDInt)).Component.Render(ctx, w)
}

func (h *Handler) NewThreadPage(w http.ResponseWriter, r *http.Request) {
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

	templ.Handler(views.Home(
		views.NewThreadForm(member.Username),
		views.NewThreadTitleGroup(),
		member.Username,
	)).Component.Render(r.Context(), w)
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

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

	body := r.FormValue("hidden_body")
	if body == "" {
		h.logger.Error("body is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ip := r.RemoteAddr
	if strings.Contains(ip, "[::1]") {
		ip = "127.0.0.1"
	}

	splitIP := strings.Split(ip, ":")
	if len(splitIP) > 1 {
		ip = splitIP[0]
	}

	threadID, err := h.threadService.NewThread(member.Username, ip, body, subject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/thread/view/%d", threadID))

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ToggleDot(w http.ResponseWriter, r *http.Request) {
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

	dotted, err := h.threadService.ToggleDot(ctx, member.ID, threadIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(commonviews.DotControl(dotted, threadIDInt)).Component.Render(ctx, w)
}

func (h *Handler) Posts(w http.ResponseWriter, r *http.Request) {
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

	thread, err := h.threadService.GetCollapsibleThreadByID(r.Context(), end-start, threadID, member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var posts []common.Post
	for _, post := range thread.Posts {
		body, err := h.imageService.PresignPostImages(r.Context(), post.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		posts = append(posts, common.Post{
			Body:       body,
			MemberName: post.MemberName,
			Date:       post.Timestamp.Format("Mon Jan 2, 2006 03:04 pm"),
			RowNumber:  post.RowNumber + 1,
			MemberID:   post.MemberID,
			ParentID:   threadID,
		})
	}

	templ.Handler(commonviews.Posts(posts, member.Viewable, member.Username, "thread")).Component.Render(r.Context(), w)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
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

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := r.FormValue("hidden_body")
	if body == "" {
		h.logger.Error("body is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	processedBody, err := h.imageService.PresignPostImages(r.Context(), body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	splitIP := strings.Split(ip, ":")
	if len(splitIP) > 1 {
		ip = splitIP[0]
	}

	_, err = h.threadService.NewPost(body, ip, member.Username, threadIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post := common.Post{
		Body:       processedBody,
		MemberName: member.Username,
		Date:       time.Now().Format("Mon Jan 2, 2006 03:04 pm"),
	}

	templ.Handler(commonviews.Post(post, idxInt, true)).Component.Render(r.Context(), w)
}

func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
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

	err = r.ParseForm()
	if err != nil {
		h.logger.Errorf("error parsing form: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := r.FormValue("hidden_body")
	if body == "" {
		h.logger.Error("body is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	processedBody, err := h.imageService.PresignPostImages(r.Context(), body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		Body:       processedBody,
		MemberName: member.Username,
		Date:       timestamp,
		Preview:    true,
	}

	templ.Handler(commonviews.Post(post, idxInt, true)).Component.Render(r.Context(), w)
}

func (h *Handler) Thread(w http.ResponseWriter, r *http.Request) {
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

	thread, err := h.threadService.GetCollapsibleThreadByID(ctx, member.Viewable, idInt, member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var posts []common.Post
	for _, post := range thread.Posts {
		body, err := h.imageService.PresignPostImages(r.Context(), post.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		posts = append(posts, common.Post{
			Body:       body,
			MemberName: post.MemberName,
			Date:       post.Timestamp.Format("Mon Jan 2, 2006 03:04 pm"),
			RowNumber:  post.RowNumber + 1,
			MemberID:   post.MemberID,
			ParentID:   thread.ID,
		})
	}

	templ.Handler(views.Home(views.Thread(commonviews.PostsPage(posts, member.Viewable, member.Username, "thread", thread.Undot), thread.ID), commonviews.PostsTitleGroup(*thread), member.Username)).Component.Render(ctx, w)
}

func (h *Handler) ThreadsHome(w http.ResponseWriter, r *http.Request) {
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

	threads, cursorsOut, err := h.threadService.ListThreads(ctx, domain.Cursors{}, 30, member.ID, domain.ThreadFilterAll)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(views.Home(views.Threads(threads, cursorsOut, member.Username), views.ThreadsTitleGroup(ElitismTitle), member.Username)).Component.Render(ctx, w)
}

func (h *Handler) Threads(w http.ResponseWriter, r *http.Request) {
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

	threads, cursors, err := h.threadService.ListThreads(ctx, cursors, 30, member.ID, domain.ThreadFilterAll)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isHx(r) {
		templ.Handler(views.Threads(threads, cursors, member.Username)).Component.Render(ctx, w)
		return
	}

	templ.Handler(views.Home(views.Threads(threads, cursors, member.Username), views.ThreadsTitleGroup(ElitismTitle), member.Username)).Component.Render(ctx, w)
}

func (h *Handler) IgnoredThreads(w http.ResponseWriter, r *http.Request) {
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

	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	memberIDInt, err := strconv.Atoi(memberID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nextCursor := r.URL.Query().Get("next")
	prevCursor := r.URL.Query().Get("prev")
	dir := r.URL.Query().Get("dir")

	if dir == "next" {
		prevCursor = ""
	} else {
		nextCursor = ""
	}

	cursors := domain.Cursors{
		Next: nextCursor,
		Prev: prevCursor,
	}

	threads, cursors, err := h.threadService.ListThreads(ctx, cursors, 50, memberIDInt, domain.ThreadFilterIgnored)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isHx(r) {
		templ.Handler(views.Threads(threads, cursors, member.Username)).Component.Render(ctx, w)
		return
	}

	viewMember, err := h.memberService.GetMemberByID(memberIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	title := fmt.Sprintf("Ignored threads: %s", viewMember.Name)

	templ.Handler(views.Home(
		views.Threads(threads, cursors, member.Username),
		views.ThreadsTitleGroup(title),
		member.Username,
	)).Component.Render(ctx, w)
}

func (h *Handler) CreatedThreads(w http.ResponseWriter, r *http.Request) {
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

	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	memberIDInt, err := strconv.Atoi(memberID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nextCursor := r.URL.Query().Get("next")
	prevCursor := r.URL.Query().Get("prev")
	dir := r.URL.Query().Get("dir")

	if dir == "next" {
		prevCursor = ""
	} else {
		nextCursor = ""
	}

	cursors := domain.Cursors{
		Next: nextCursor,
		Prev: prevCursor,
	}

	threads, cursors, err := h.threadService.ListThreads(ctx, cursors, 50, memberIDInt, domain.ThreadFilterCreated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isHx(r) {
		templ.Handler(views.Threads(threads, cursors, member.Username)).Component.Render(ctx, w)
		return
	}

	viewMember, err := h.memberService.GetMemberByID(memberIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	title := fmt.Sprintf("Created threads: %s", viewMember.Name)

	templ.Handler(views.Home(views.Threads(
		threads,
		cursors,
		member.Username,
	), views.ThreadsTitleGroup(title), member.Username)).Component.Render(ctx, w)
}

func (h *Handler) ParticipatedThreads(w http.ResponseWriter, r *http.Request) {
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

	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	memberIDInt, err := strconv.Atoi(memberID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nextCursor := r.URL.Query().Get("next")
	prevCursor := r.URL.Query().Get("prev")
	dir := r.URL.Query().Get("dir")

	if dir == "next" {
		prevCursor = ""
	} else {
		nextCursor = ""
	}

	cursors := domain.Cursors{
		Next: nextCursor,
		Prev: prevCursor,
	}

	threads, cursors, err := h.threadService.ListThreads(ctx, cursors, 50, memberIDInt, domain.ThreadFilterParticipated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isHx(r) {
		templ.Handler(views.Threads(threads, cursors, member.Username)).Component.Render(ctx, w)
		return
	}

	viewMember, err := h.memberService.GetMemberByID(memberIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	title := fmt.Sprintf("Participated threads: %s", viewMember.Name)

	templ.Handler(views.Home(views.Threads(
		threads,
		cursors,
		member.Username,
	), views.ThreadsTitleGroup(title), member.Username)).Component.Render(ctx, w)
}

func (h *Handler) FavoritedThreads(w http.ResponseWriter, r *http.Request) {
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

	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	memberIDInt, err := strconv.Atoi(memberID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nextCursor := r.URL.Query().Get("next")
	prevCursor := r.URL.Query().Get("prev")
	dir := r.URL.Query().Get("dir")

	if dir == "next" {
		prevCursor = ""
	} else {
		nextCursor = ""
	}

	cursors := domain.Cursors{
		Next: nextCursor,
		Prev: prevCursor,
	}

	threads, cursors, err := h.threadService.ListThreads(ctx, cursors, 50, memberIDInt, domain.ThreadFilterFavorites)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isHx(r) {
		templ.Handler(views.Threads(threads, cursors, member.Username)).Component.Render(ctx, w)
		return
	}

	viewMember, err := h.memberService.GetMemberByID(memberIDInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	title := fmt.Sprintf("Favorited threads: %s", viewMember.Name)

	templ.Handler(views.Home(views.Threads(
		threads,
		cursors,
		member.Username,
	), views.ThreadsTitleGroup(title), member.Username)).Component.Render(ctx, w)
}

func (h *Handler) Uploader(w http.ResponseWriter, r *http.Request) {
	templ.Handler(views.Home(commonviews.Uploader(), views.NewThreadTitleGroup(), "gofreescout")).Component.Render(r.Context(), w)
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Errorw("failed to parse form", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		h.logger.Errorw("failed to get file", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		h.logger.Errorw("failed to read file", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	key, url, err := h.imageService.UploadImage(ctx, imageBytes)
	if err != nil {
		var formatError gberrors.UnsupportedImageFormat
		if errors.As(err, &formatError) {
			w.WriteHeader(http.StatusBadRequest)
			templ.Handler(commonviews.PostError("Unsupported image format")).Component.Render(ctx, w)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(commonviews.InitialImage(url, key)).Component.Render(ctx, w)
}

func isHx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
