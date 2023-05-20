package thread

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"goBoard/internal/core/ports"
	"strconv"
)

type TemplateHandler struct {
	threadService ports.ThreadService
}

func NewTemplateHandler(threadService ports.ThreadService) *TemplateHandler {
	return &TemplateHandler{threadService}
}

func (h *TemplateHandler) Register(e *echo.Echo) {
	e.GET("/threads/all", h.ListThreads)
	e.GET("/thread/:id", h.ListPostsForThread)
	e.GET("/ping", h.Ping)
}

func (h *TemplateHandler) ListThreads(c echo.Context) error {
	threads, err := h.threadService.ListThreads(10, 0)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.Render(200, "main", threads)
}

func (h *TemplateHandler) Ping(c echo.Context) error {
	return c.Render(200, "ping", nil)
}

func (h *TemplateHandler) ListPostsForThread(c echo.Context) error {
	threadID := c.Param("id")

	idAsInt, err := strconv.Atoi(threadID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	posts, err := h.threadService.GetThreadByID(10, 0, idAsInt)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	log.Infof("posts: %+v", posts)

	return c.Render(200, "posts", posts)
}
