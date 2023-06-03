package member

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"goBoard/helpers/auth"
	"goBoard/internal/core/ports"

	"github.com/labstack/echo/v4"
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

func (h *TemplateHandler) Register(echo *echo.Echo) {
	e := echo.Group("/member")

	e.Use(echojwt.WithConfig(echojwt.Config{
		//NewClaimsFunc: auth.GetJWTClaims,
		SigningKey:   []byte(auth.GetJWTSecret()),
		TokenLookup:  "cookie:access-token", // "<source>:<name>"
		ErrorHandler: auth.JWTErrorChecker,
	}))

	e.GET("/view/:username", h.GetMemberByUsername)
}

func (h *TemplateHandler) GetMemberByUsername(c echo.Context) error {
	username := c.Param("username")
	siteContext, err := h.memberService.GetMemberByUsername(username)
	if err != nil {
		c.String(500, err.Error())
		return err
	}
	siteContext.PageName = "member"
	return c.Render(200, "member", siteContext)
}

type GenericResponse struct {
	Message string `json:"message"`
}
