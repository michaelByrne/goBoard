package message

import (
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/handlersold/member"
	"strconv"
	"strings"
)

type Handler struct {
	memberService  ports.MemberService
	messageService ports.MessageService
}

func NewHandler(memberService ports.MemberService, messageService ports.MessageService) *Handler {
	return &Handler{
		messageService: messageService,
		memberService:  memberService,
	}
}

func (h *Handler) Register(e *echo.Echo) {
	e.POST("message/addmember", h.AddMember)
	e.POST("message/create", h.CreateMessage)
	e.POST("message/reply", h.MessageReply)
}

func (h *Handler) AddMember(c echo.Context) error {
	namesString := c.FormValue("names")

	names := strings.Split(namesString, ",")
	var cleanNames []string
	for _, name := range names {
		cleanNames = append(cleanNames, strings.TrimSpace(name))
	}

	validMembers, err := h.memberService.ValidateMembers(cleanNames)
	if err != nil {
		c.String(500, err.Error())
	}

	membersOut := &member.Members{}
	membersOut.FromDomain(validMembers)

	return c.JSON(200, membersOut)
}

func (h *Handler) CreateMessage(c echo.Context) error {
	subject := c.FormValue("subject")
	body := c.FormValue("body")
	membersString := c.FormValue("message_members")
	memberID := c.FormValue("member_id")

	memberIDAsInt, err := strconv.Atoi(memberID)
	if err != nil {
		return c.String(500, err.Error())
	}

	memberIP := c.RealIP()

	members := strings.Split(membersString, ",")
	var cleanMembersAsInts []int
	for _, m := range members {
		cleanMember := strings.TrimSpace(m)
		cleanMemberAsInt, err := strconv.Atoi(cleanMember)
		if err != nil {
			return c.String(500, err.Error())
		}

		cleanMembersAsInts = append(cleanMembersAsInts, cleanMemberAsInt)
	}

	messageID, err := h.messageService.SendMessage(subject, body, memberIP, memberIDAsInt, cleanMembersAsInts)
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.JSON(200, NewMessageResponse{
		ID: messageID,
	})
}

func (h *Handler) MessageReply(c echo.Context) error {
	messageID := c.FormValue("message_id")

	messageIDAsInt, err := strconv.Atoi(messageID)
	if err != nil {
		return c.String(500, err.Error())
	}

	body := c.FormValue("body")
	memberName := c.FormValue("member_name")

	postID, err := h.messageService.NewPost(body, c.RealIP(), memberName, messageIDAsInt)
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.JSON(200, NewPostResponse{
		ID: postID,
	})
}

type MembersRequest struct {
	Names []string `json:"names"`
}

type NewMessageResponse struct {
	ID int `json:"id"`
}

type NewPostResponse struct {
	ID int `json:"post_id"`
}
