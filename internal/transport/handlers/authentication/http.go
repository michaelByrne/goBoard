package authentication

import (
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/ports"
)

type HTTPHandler struct {
	authenticationService ports.AuthenticationService
}

func NewHTTPHandler(authenticationService ports.AuthenticationService) *HTTPHandler {
	return &HTTPHandler{
		authenticationService: authenticationService,
	}
}

func (h *HTTPHandler) Register(e *echo.Echo) {
	e.POST("/login", h.Login)
	//e.POST("/logout", h.Logout)
}

func (h *HTTPHandler) Login(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("pass")

	memberID, err := h.authenticationService.Authenticate(c.Request().Context(), username, password)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	if memberID == 0 {
		c.String(401, "Unauthorized")
		return nil
	}

	return c.Redirect(302, "/")
}
