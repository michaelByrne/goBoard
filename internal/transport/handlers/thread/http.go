package thread

import (
	"errors"
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/middlewares/session"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	threadService      ports.ThreadService
	defaultThreadLimit int
}

func NewHandler(threadService ports.ThreadService, defaultThreadLimit int) *Handler {
	return &Handler{
		threadService:      threadService,
		defaultThreadLimit: defaultThreadLimit,
	}
}

func (h *Handler) Register(e *echo.Echo) {
	e.GET("/threads/:id", h.GetThreadByID)
	e.POST("/thread/create", h.CreateThread)
	e.POST("/thread/reply", h.ThreadReply)
	//e.GET("/threads/cursor", h.GetThreadsWithCursor)
	e.GET("/threads/home", h.ListThreads)
	e.POST("thread/undot/:id", h.UndotThread)
	e.POST("thread/ignore", h.ToggleIgnore)

}

func (h *Handler) ThreadReply(c echo.Context) error {
	values, err := c.FormParams()
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	threadID := values.Get("thread_id")
	threadIDAsInt, err := strconv.Atoi(threadID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	author := values.Get("member_name")
	body := values.Get("body")

	ip := c.RealIP()

	postID, err := h.threadService.NewPost(body, ip, author, threadIDAsInt)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.JSON(200, NewPostResponse{
		PostID: postID,
	})
}

func (h *Handler) GetThreadByID(ctx echo.Context) error {
	sess, err := session.Get("member", ctx)
	if err != nil {
		ctx.String(500, err.Error())
		return err
	}

	memberID, ok := sess.Values["id"]
	if !ok {
		ctx.String(500, "member id not found in session")
		return errors.New("member id not found in session")
	}

	memberIDAsInt, ok := memberID.(int)
	if !ok {
		ctx.String(500, "member id not an int")
		return errors.New("member id not an int")
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
		return err
	}

	thread, err := h.threadService.GetThreadByID(100, 100, id, memberIDAsInt)
	if err != nil {
		ctx.JSON(500, ErrorResponse{Message: err.Error()})
		return err
	}

	threadOut := &Thread{}
	threadOut.FromDomain(*thread)

	return ctx.JSON(200, threadOut)
}

func (h *Handler) CreateThread(c echo.Context) error {
	body := c.FormValue("body")
	subject := c.FormValue("subject")

	ip := c.RealIP()

	author := c.FormValue("member")

	threadID, err := h.threadService.NewThread(author, ip, body, subject)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.JSON(200, NewThreadResponse{
		ThreadID: threadID,
	})
}

func (h *Handler) ListThreads(c echo.Context) error {
	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	memberID, ok := sess.Values["id"]
	if !ok {
		c.String(500, "member id not found")
		return err
	}

	memberIDAsInt, ok := memberID.(int)
	if !ok {
		c.String(500, "member id not int")
		return err
	}

	cursor := c.QueryParams().Get("cursor")
	var cursorAsTime time.Time
	if cursor == "" || cursor == "null" {
		cursorAsTime = time.Date(9999, 1, 1, 1, 1, 1, 1, time.UTC)
	} else {
		cursorAsTime, err = time.Parse(time.RFC3339, cursor)
		if err != nil {
			c.String(500, err.Error())
			return err
		}
	}

	reverse := c.QueryParams().Get("reverse")
	var reverseAsBool bool
	if reverse == "" {
		reverseAsBool = false
	} else {
		reverseAsBool, err = strconv.ParseBool(reverse)
		if err != nil {
			c.String(500, err.Error())
			return err
		}
	}

	if reverseAsBool {
		siteContext, err := h.threadService.GetThreadsWithCursorReverse(h.defaultThreadLimit, &cursorAsTime, memberIDAsInt)
		if err != nil {
			c.String(500, err.Error())
			return err
		}

		siteContext.PageName = "main"
		return c.JSON(200, siteContext)
	}

	siteContext, err := h.threadService.GetThreadsWithCursorForward(h.defaultThreadLimit, false, &cursorAsTime, memberIDAsInt)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	siteContext.PageName = "main"
	return c.JSON(200, siteContext)
}

func (h *Handler) UndotThread(c echo.Context) error {
	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	memberID, ok := sess.Values["id"]
	if !ok {
		c.String(500, "member id not found")
		return err
	}

	memberIDAsInt, ok := memberID.(int)
	if !ok {
		c.String(500, "member id not int")
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	err = h.threadService.UndotThread(c.Request().Context(), memberIDAsInt, id)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.String(200, "success")
}

func (h *Handler) ToggleIgnore(c echo.Context) error {
	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	memberID, ok := sess.Values["id"]
	if !ok {
		c.String(500, "member id not found")
		return err
	}

	memberIDAsInt, ok := memberID.(int)
	if !ok {
		c.String(500, "member id not int")
		return err
	}

	var ignoreRequest IgnoreRequest
	err = c.Bind(&ignoreRequest)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	err = h.threadService.ToggleIgnore(c.Request().Context(), memberIDAsInt, ignoreRequest.ThreadID, ignoreRequest.Ignore)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.String(200, "success")
}

type IgnoreRequest struct {
	ThreadID int  `json:"id"`
	Ignore   bool `json:"ignore"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type NewPostResponse struct {
	PostID int `json:"post_id"`
}

type NewThreadResponse struct {
	ThreadID int `json:"thread_id"`
}
