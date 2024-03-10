package message

import (
	"errors"
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/middlewares/session"
	"strconv"
	"time"
)

type TemplateHandler struct {
	messageService ports.MessageService
}

func NewTemplateHandler(messageService ports.MessageService) *TemplateHandler {
	return &TemplateHandler{
		messageService: messageService,
	}
}

func (h *TemplateHandler) Register(echo *echo.Echo) {
	e := echo.Group("/message")

	e.GET("/create", h.CreateMessage)
	e.GET("/list/:memberID", h.ListMessages)
	e.GET("/view/:id", h.ViewMessage)
	e.POST("/previewpost/:position", h.PreviewPost)
	e.GET("/post/:id/:position", h.Post)
}

func (h *TemplateHandler) CreateMessage(c echo.Context) error {
	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	memberName, ok := sess.Values["name"]
	if !ok {
		c.String(500, "no member name in session")
		return errors.New("no member name in session")
	}

	memberNameAsStr, ok := memberName.(string)
	if !ok {
		c.String(500, "member name not a string")
		return errors.New("member name not a string")
	}

	siteContext := domain.SiteContext{
		Username: memberNameAsStr,
	}

	return c.Render(200, "new-message", siteContext)
}

func (h *TemplateHandler) ListMessages(c echo.Context) error {
	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	memberName, ok := sess.Values["name"]
	if !ok {
		c.String(500, "no member name in session")
		return err
	}

	memberNameAsStr, ok := memberName.(string)
	if !ok {
		c.String(500, "member name not a string")
		return err
	}

	memberID := c.Param("memberID")

	cursor := c.QueryParams().Get("cursor")
	var cursorAsTime time.Time
	if cursor == "" {
		cursorAsTime = time.Date(2999, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		cursorAsTime, err = time.Parse(time.RFC3339, cursor)
		if err != nil {
			c.String(500, err.Error())
			return err
		}
	}

	reverse := c.QueryParams().Get("reverse")
	reverseAsBool, err := strconv.ParseBool(reverse)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	memberIDAsInt, err := strconv.Atoi(memberID)
	if err != nil {
		return c.String(500, err.Error())
	}

	messages, err := h.messageService.GetMessagesWithCursor(memberIDAsInt, reverseAsBool, &cursorAsTime)
	if err != nil {
		return c.String(500, err.Error())
	}

	siteContext := domain.SiteContext{
		Messages: messages,
		Username: memberNameAsStr,
	}

	return c.Render(200, "messages", siteContext)
}

func (h *TemplateHandler) ViewMessage(c echo.Context) error {
	id := c.Param("id")

	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return c.String(500, err.Error())
	}

	message, err := h.messageService.GetMessageByID(idAsInt, 1)

	return c.Render(200, "message", message)
}

func (h *TemplateHandler) PreviewPost(c echo.Context) error {
	position := c.Param("position")

	positionAsInt, err := strconv.Atoi(position)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	values, err := c.FormParams()
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	body := values.Get("body")
	messageID := values.Get("message_id")
	author := values.Get("member_name")

	messageIDAsInt, err := strconv.Atoi(messageID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	now := time.Now()

	post := &domain.MessagePost{
		ParentID:   messageIDAsInt,
		Body:       body,
		MemberName: author,
		Timestamp:  &now,
		Position:   positionAsInt + 1,
	}

	return c.Render(200, "post", post)
}

func (h *TemplateHandler) Post(c echo.Context) error {
	postID := c.Param("id")
	idAsInt, err := strconv.Atoi(postID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	post, err := h.messageService.GetMessagePostByID(idAsInt)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.Render(200, "post", post)
}
