package thread

import (
	"goBoard/internal/core/ports"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	threadService ports.ThreadService
}

func NewHandler(threadService ports.ThreadService) *Handler {
	return &Handler{threadService}
}

func (h *Handler) Register(e *echo.Echo) {
	e.GET("/threads/all", h.ListThreads)
	e.GET("/threads/:id", h.GetThreadByID)
	e.GET("/thread/create", h.NewThread)
	e.POST("/thread/create", h.CreateThread)
	e.POST("/thread/reply", h.ThreadReply)
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

func (h *Handler) ListThreads(ctx echo.Context) error {
	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
		return err
	}

	offset, err := strconv.Atoi(ctx.QueryParam("offset"))
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
		return err
	}

	threadList, err := h.threadService.ListThreads(limit, offset)
	if err != nil {
		ctx.JSON(500, ErrorResponse{Message: err.Error()})
		return err
	}

	return ctx.JSON(200, threadList)
}

func (h *Handler) GetThreadByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
		return err
	}

	thread, err := h.threadService.GetThreadByID(100, 100, id)
	if err != nil {
		ctx.JSON(500, ErrorResponse{Message: err.Error()})
		return err
	}

	threadOut := &Thread{}
	threadOut.FromDomain(*thread)

	return ctx.JSON(200, threadOut)
}

//	func (h *Handler) SavePost(ctx echo.Context) error {
//		post := &Post{}
//		err := ctx.Bind(post)
//		if err != nil {
//			ctx.JSON(400, ErrorResponse{Message: err.Error()})
//			return err
//		}
//
//		id, err := h.threadService.Save(post.ToDomain())
//		if err != nil {
//			ctx.JSON(500, ErrorResponse{Message: err.Error()})
//			return err
//		}
//
//		return ctx.JSON(200, ID{id})
//	}

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

//
//func (h *Handler) NewThread(c echo.Context) error {
//	thread := &Thread{}
//	err := c.Bind(thread)
//	if err != nil {
//		c.JSON(400, ErrorResponse{Message: err.Error()})
//		return err
//	}
//
//	id, err := h.threadService.NewThread(strconv.Itoa(thread.MemberID), thread.MemberIP, thread.FirstPostText, thread.Subject)
//	if err != nil {
//		c.JSON(500, ErrorResponse{Message: err.Error()})
//		return err
//	}
//
//	return c.JSON(200, ID{ID: id})
//}

func (h *Handler) NewThread(c echo.Context) error {
	thread := &Thread{}
	err := c.Bind(thread)
	if err != nil {
		c.JSON(400, ErrorResponse{Message: err.Error()})
		return err
	}

	id, err := h.threadService.NewThread(strconv.Itoa(thread.MemberID), thread.MemberIP, thread.FirstPostText, thread.Subject)
	if err != nil {
		c.JSON(500, ErrorResponse{Message: err.Error()})
		return err
	}

	return c.JSON(200, ID{ID: id})
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
