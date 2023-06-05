package message

import (
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
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
	return c.Render(200, "new-message", nil)
}

func (h *TemplateHandler) ListMessages(c echo.Context) error {
	memberID := c.Param("memberID")

	cursor := c.QueryParams().Get("cursor")
	var cursorAsTime time.Time
	var err error
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

	return c.Render(200, "messages", messages)
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
