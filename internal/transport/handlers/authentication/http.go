package authentication

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"goBoard/helpers/auth"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/middlewares/session"
	"time"
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
	e.GET("/logout/:name", h.Logout)
}

func (h *HTTPHandler) Login(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("pass")

	member, err := h.authenticationService.Authenticate(c.Request().Context(), username, password)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	if member == nil {
		c.String(401, "Unauthorized")
		return nil
	}

	err = auth.GenerateTokensAndSetCookies(&domain.Member{
		ID:   member.ID,
		Name: username,
		Pass: password,
	}, c, time.Hour)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	sess.Values["id"] = member.ID
	sess.Values["name"] = username
	sess.Values["admin"] = member.IsAdmin

	for key, pref := range member.Prefs {
		sess.Values[key] = pref.Value
	}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.Redirect(302, "/")
}

func (h *HTTPHandler) Logout(c echo.Context) error {
	username := c.Param("name")

	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	err = auth.GenerateTokensAndSetCookies(&domain.Member{Name: username}, c, -1)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.Redirect(302, "/")
}
