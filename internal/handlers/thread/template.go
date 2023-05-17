package thread

import (
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/ports"
)

type TemplateHandler struct {
	threadService ports.ThreadService
}

func NewTemplateHandler(threadService ports.ThreadService) *TemplateHandler {
	return &TemplateHandler{threadService}
}

func (h *TemplateHandler) Register(e *echo.Echo) {
	e.GET("/threads/all", h.ListThreads)
	e.GET("/ping", h.Ping)
}

func (h *TemplateHandler) ListThreads(c echo.Context) error {
	threads, err := h.threadService.ListThreads(10, 0)
	if err != nil {
		c.String(500, err.Error())
	}

	return c.Render(200, "main", threads)
}

func (h *TemplateHandler) Ping(c echo.Context) error {
	return c.Render(200, "ping", nil)
}
