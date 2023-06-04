package member

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"goBoard/helpers/auth"
	"goBoard/internal/core/ports"
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

	e.Use(echojwt.WithConfig(echojwt.Config{
		//NewClaimsFunc: auth.GetJWTClaims,
		SigningKey:   []byte(auth.GetJWTSecret()),
		TokenLookup:  "cookie:access-token", // "<source>:<name>"
		ErrorHandler: auth.JWTErrorChecker,
	}))

	e.POST("/save", h.SaveMember)
	e.GET("/:id", h.GetMemberByID)
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

type ErrorResponse struct {
	Message string `json:"message"`
}
