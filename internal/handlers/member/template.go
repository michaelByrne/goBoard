package member

import (
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/ports"
)

type TemplateHandler struct {
	threadService ports.ThreadService
	memberService ports.MemberService
}

func NewTemplateHandler(threadService ports.ThreadService, memberService ports.MemberService) *TemplateHandler {
	return &TemplateHandler{
		threadService: threadService,
		memberService: memberService,
	}
}

func (h *TemplateHandler) Register(e *echo.Echo) {
	e.GET("/member/view/:username", h.GetMemberByUsername)
}


func (h *TemplateHandler) GetMemberByUsername(c echo.Context) error {
	username := c.Param("username")
	member, err := h.memberService.GetMemberByUsername(username)
	if err != nil {
		c.String(500, err.Error())
		return err
	}
	return c.Render(200, "member", member)
}

type GenericResponse struct {
	Message string `json:"message"`
}
