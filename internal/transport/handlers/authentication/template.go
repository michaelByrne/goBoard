package authentication

import (
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/ports"
)

type TemplateHandler struct {
	authenticationService ports.AuthenticationService
}

func NewTemplateHandler(authenticationService ports.AuthenticationService) *TemplateHandler {
	return &TemplateHandler{
		authenticationService: authenticationService,
	}
}

func (h *TemplateHandler) Register(e *echo.Echo) {
	e.GET("/login", h.Login)
	//e.POST("/logout", h.Logout)
}

func (h *TemplateHandler) Login(c echo.Context) error {
	return c.Render(200, "login", nil)
}
