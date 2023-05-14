package thread

import (
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/ports"
	"strconv"
)

type Handler struct {
	threadService ports.ThreadService
}

func NewHandler(threadService ports.ThreadService) *Handler {
	return &Handler{threadService}
}

func (h *Handler) Register(e *echo.Echo) {
	e.GET("/threads/:limit", h.ListThreads)
	e.GET("/threads/:id", h.GetThreadByID)
	e.POST("/threads/posts", h.SavePost)
	e.POST("/threads", h.NewThread)
}

func (h *Handler) ListThreads(ctx echo.Context) error {
	limit, err := strconv.Atoi(ctx.Param("limit"))
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
	}

	threads, err := h.threadService.ListThreads(limit)
	if err != nil {
		ctx.JSON(500, ErrorResponse{Message: err.Error()})
	}

	threadsOut := &Threads{}
	threadsOut.FromDomain(threads)

	return ctx.JSON(200, threadsOut)
}

func (h *Handler) GetThreadByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
	}

	thread, err := h.threadService.GetThreadByID(id)
	if err != nil {
		ctx.JSON(500, ErrorResponse{Message: err.Error()})
	}

	threadOut := &Thread{}
	threadOut.FromDomain(*thread)

	return ctx.JSON(200, threadOut)
}

func (h *Handler) SavePost(ctx echo.Context) error {
	post := &Post{}
	err := ctx.Bind(post)
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
	}

	err = h.threadService.Save(post.ToDomain())
	if err != nil {
		ctx.JSON(500, ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(200, nil)
}

func (h *Handler) NewThread(c echo.Context) error {
	thread := &Thread{}
	err := c.Bind(thread)
	if err != nil {
		c.JSON(400, ErrorResponse{Message: err.Error()})
	}

	err = h.threadService.NewThread(thread.ToDomain())
	if err != nil {
		c.JSON(500, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(200, nil)
}

type ErrorResponse struct {
	Message string `json:"message"`
}
