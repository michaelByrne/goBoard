package member

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"goBoard/helpers/auth"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/middlewares/session"
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
	e.GET("/edit", h.EditMember)
}

func (h *TemplateHandler) GetMemberByUsername(c echo.Context) error {
	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	siteContext := domain.SiteContext{
		Session: sess,
	}

	username := c.Param("username")
	member, err := h.memberService.GetMemberByUsername(username)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	siteContext.Member = *member

	return c.Render(200, "member", siteContext)
}

func (h *TemplateHandler) EditMember(c echo.Context) error {
	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	memberID := sess.Values["id"].(int)
	memberName := sess.Values["name"].(string)

	siteContext := domain.SiteContext{
		Session:  sess,
		Username: memberName,
	}

	prefs, err := h.memberService.GetMergedPrefs(c.Request().Context(), memberID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	siteContext.Prefs = prefs

	member, err := h.memberService.GetMemberByID(memberID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	siteContext.Member = *member

	return c.Render(200, "member-edit", siteContext)
}

type GenericResponse struct {
	Message string `json:"message"`
}
