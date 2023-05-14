package member

import (
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/ports"
	"strconv"
)

type Handler struct {
	memberService ports.MemberService
}

func NewHandler(memberService ports.MemberService) *Handler {
	return &Handler{memberService}
}

func (h *Handler) Register(e *echo.Echo) {
	e.POST("/members", h.SaveMember)
	e.GET("/members/:id", h.GetMemberByID)
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

type ErrorResponse struct {
	Message string `json:"message"`
}
