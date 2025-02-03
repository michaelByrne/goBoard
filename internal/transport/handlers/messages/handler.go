package messages

import (
	"fmt"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/handlers/common"
	commonviews "goBoard/internal/transport/handlers/common/views"
	"goBoard/internal/transport/handlers/messages/views"
	threadviews "goBoard/internal/transport/handlers/threads/views"
	"goBoard/internal/transport/middlewares/session"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	messageService ports.MessageService
	memberService  ports.MemberService
	imageService   ports.ImageService

	verifyMiddleware func(next http.Handler) http.Handler

	logger *zap.SugaredLogger
}

func NewHandler(
	messageService ports.MessageService,
	memberService ports.MemberService,
	imageService ports.ImageService,
	verifyMiddleware func(next http.Handler) http.Handler,
	logger *zap.SugaredLogger,
) *Handler {
	return &Handler{
		messageService:   messageService,
		memberService:    memberService,
		imageService:     imageService,
		verifyMiddleware: verifyMiddleware,
		logger:           logger,
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(h.verifyMiddleware)

		r.Get("/message/create", h.NewMessagePage)
		r.Post("/message/create", h.CreateMessage)
		r.Post("/message/post", h.Post)
		r.Get("/view/message/{id}", h.ViewMessage)
		r.Post("/message/preview", h.Preview)
		r.Get("/message/view/{id}", h.Message)
		r.Get("/message/posts", h.Posts)
		r.Get("/message/list", h.Messages)
		r.Get("/message/delete/{id}", h.DeleteMessage)
		r.Get("/message/counts", h.MessageCounts)
	})
}

func (h *Handler) MessageCounts(w http.ResponseWriter, r *http.Request) {
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

	counts, err := h.messageService.GetNewMessageCounts(ctx, member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(views.MessageCounts(*counts)).Component.Render(ctx, w)
}

func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
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

	messageIDStr := chi.URLParam(r, "id")
	if messageIDStr == "" {
		h.logger.Errorf("message id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		h.logger.Errorf("error converting message id to int: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.messageService.DeleteMessage(ctx, member.ID, messageID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", "/message/list")

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
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

	templ.Handler(commonviews.Post(post, idxInt, true)).Component.Render(ctx, w)
}

func (h *Handler) ViewMessage(w http.ResponseWriter, r *http.Request) {
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

	messageIDStr := chi.URLParam(r, "id")
	if messageIDStr == "" {
		h.logger.Errorf("message id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageID, error := strconv.Atoi(messageIDStr)
	if error != nil {
		h.logger.Errorf("error converting message id to int: %v", error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.messageService.ViewMessage(ctx, member.ID, messageID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Messages(w http.ResponseWriter, r *http.Request) {
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

	if dir == "next" {
		prevCursor = ""
	} else {
		nextCursor = ""
	}

	cursors := domain.Cursors{
		Next: nextCursor,
		Prev: prevCursor,
	}

	messages, cursors, err := h.messageService.ListMessages(r.Context(), cursors, 10, member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isHx(r) {
		templ.Handler(views.Messages(messages, cursors, member.Username)).Component.Render(r.Context(), w)
		return
	}

	templ.Handler(threadviews.Home(views.Messages(messages, cursors, member.Username), views.MessagesTitleGroup("elitism. secrecy. tradition."), member.Username)).Component.Render(r.Context(), w)
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
		h.logger.Errorf("error parsing form: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
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

	_, err = h.messageService.NewPost(processedBody, ip, member.Username, threadIDInt)
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

func (h *Handler) Message(w http.ResponseWriter, r *http.Request) {
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

	messageIDStr := chi.URLParam(r, "id")
	if messageIDStr == "" {
		h.logger.Errorf("message id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageID, ok := strconv.Atoi(messageIDStr)
	if ok != nil {
		h.logger.Errorf("error converting message id to int: %v", ok)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message, err := h.messageService.GetCollapsibleMessageByID(r.Context(), 5, messageID, member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var posts []common.Post
	for _, post := range message.Posts {
		body, err := h.imageService.PresignPostImages(r.Context(), post.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		posts = append(posts, common.Post{
			Body:       body,
			MemberName: post.MemberName,
			Date:       post.Timestamp.Format("Mon Jan 2, 2006 03:04 pm"),
			RowNumber:  post.Position + 1,
			MemberID:   post.MemberID,
			ParentID:   message.ID,
		})
	}

	templ.Handler(threadviews.Home(commonviews.PostsPage(posts, 5, member.Username, "message", false), views.MessageTitleGroup(*message), member.Username)).Component.Render(r.Context(), w)
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subject := r.FormValue("subject")
	if subject == "" {
		h.logger.Errorf("subject is empty")
		w.WriteHeader(http.StatusBadRequest)
		templ.Handler(views.MessageError("Subject is empty")).Component.Render(r.Context(), w)
		return
	}

	body := r.FormValue("hidden_body")
	if body == "" {
		h.logger.Errorf("body is empty")
		w.WriteHeader(http.StatusBadRequest)
		templ.Handler(views.MessageError("Body is empty")).Component.Render(r.Context(), w)
		return
	}

	recipientsStr := r.FormValue("message_members")
	if recipientsStr == "" {
		h.logger.Errorf("recipients is empty")
		w.WriteHeader(http.StatusBadRequest)
		templ.Handler(views.MessageError("Recipients is empty")).Component.Render(r.Context(), w)
		return
	}

	recipients := strings.Split(recipientsStr, ",")

	var recipientIDs []int
	for _, recipient := range recipients {
		member, err := h.memberService.GetMemberByUsername(recipient)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		recipientIDs = append(recipientIDs, member.ID)
	}

	memberIP := r.RemoteAddr

	if strings.Contains(memberIP, "[::1]") {
		memberIP = "127.0.0.1"
	}

	splitIP := strings.Split(memberIP, ":")
	if len(splitIP) > 1 {
		memberIP = splitIP[0]
	}

	messageID, err := h.messageService.SendMessage(subject, body, memberIP, member.ID, recipientIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/message/view/%d", messageID))

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) NewMessagePage(w http.ResponseWriter, r *http.Request) {
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

	templ.Handler(threadviews.Home(views.NewMessageForm(member.Username), views.NewMessageTitleGroup(), member.Username)).Component.Render(r.Context(), w)
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

	message, err := h.messageService.GetCollapsibleMessageByID(r.Context(), end-start, threadID, member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	posts := common.MessageToPosts(*message)

	templ.Handler(commonviews.Posts(posts, 5, member.Username, "message")).Component.Render(r.Context(), w)
}

func isHx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
