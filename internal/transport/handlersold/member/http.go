package member

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/middlewares/session"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	memberService ports.MemberService
}

func NewHandler(memberService ports.MemberService) *Handler {
	return &Handler{memberService}
}

func (h *Handler) Register(echo *echo.Echo) {
	e := echo.Group("/member")

	e.POST("/save", h.SaveMember)
	e.GET("/:id", h.GetMemberByID)
	e.POST("/edit", h.EditMember)
	//e.GET("/member/view/:username", h.GetMemberByUsername)
}

func (h *Handler) SaveMember(ctx echo.Context) error {
	member := &Member{}
	err := ctx.Bind(member)
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
		return err
	}

	id, err := h.memberService.Save(member.ToDomain())
	if err != nil {
		ctx.JSON(500, ErrorResponse{Message: err.Error()})
		return err
	}

	return ctx.JSON(200, ID{id})
}

func (h *Handler) GetMemberByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, ErrorResponse{Message: err.Error()})
		return err
	}

	member, err := h.memberService.GetMemberByID(id)
	if err != nil {
		ctx.JSON(500, ErrorResponse{Message: err.Error()})
		return err
	}

	memberOut := &Member{}
	memberOut.FromDomain(*member)

	return ctx.JSON(200, memberOut)
}

func (h *Handler) GetMemberByUsername(ctx echo.Context) error {
	username := ctx.Param("username")

	siteContext, err := h.memberService.GetMemberByUsername(username)
	if err != nil {
		ctx.JSON(500, ErrorResponse{Message: err.Error()})
		return err
	}

	return ctx.JSON(200, siteContext)
}

func (h *Handler) EditMember(c echo.Context) error {
	sess, err := session.Get("member", c)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	memberID := sess.Values["id"].(int)

	values, err := c.FormParams()
	if err != nil {
		return err
	}

	var updatedCount int
	prefs := make(domain.MemberPrefs)
	for k, v := range values {
		if len(v) == 0 {
			continue
		}

		if k == "username" {
			continue
		}

		if k == "postal" {
			member, err := h.memberService.GetMemberByID(memberID)
			if err != nil {
				c.String(500, err.Error())
				return err
			}

			member.PostalCode = v[0]

			err = h.memberService.UpdateMember(c.Request().Context(), *member)
			if err != nil {
				c.String(500, err.Error())
				return err
			}

			continue
		}

		if v[0] != "" {
			var value string
			if len(v) == 2 {
				value = v[1]
			} else {
				value = v[0]
			}
			prefs[k] = domain.MemberPref{
				Value: value,
			}
			updatedCount++
		}
	}

	err = h.memberService.UpdatePrefs(c.Request().Context(), memberID, prefs)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.JSON(200, SuccessMessage{Message: "Updated " + strconv.Itoa(updatedCount) + " preferences"})
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessMessage struct {
	Message string `json:"message"`
}
