package thread

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"strconv"
	"time"

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

func (h *TemplateHandler) Register(e *echo.Echo) {
	e.GET("/", h.ListFirstPageThreads)
	e.GET("/thread/list", h.ListThreads)
	e.GET("/thread/view/:id", h.ListPostsForThread)
	e.GET("/post/:id/:position", h.Post)
	e.GET("/ping", h.Ping)
	e.POST("/thread/reply", h.ThreadReply)
	e.POST("/thread/create", h.CreateThread)
	e.GET("/thread/create", h.NewThread)
	e.POST("/thread/previewpost/:position", h.PreviewPost)
}

func (h *TemplateHandler) ListFirstPageThreads(c echo.Context) error {
	threadListLength := 3
	siteContext, err := h.threadService.GetThreadsWithCursor(threadListLength, true, nil)
	if err != nil {
		c.String(500, err.Error())
		return err
	}
	siteContext.ThreadPage.PageNum = 0
	siteContext.PageName = "main"
	return c.Render(200, "main", siteContext)
}

func (h *TemplateHandler) ListThreads(c echo.Context) error {
	cursor := c.QueryParams().Get("cursor")
	limit, err := strconv.Atoi(c.QueryParams().Get("limit"))
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	cursorAsTime, err := time.Parse(time.RFC3339, cursor)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	reverse := c.QueryParams().Get("reverse")
	reverseAsBool, err := strconv.ParseBool(reverse)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	if reverseAsBool {
		threads, err := h.threadService.GetThreadsWithCursorReverse(limit, &cursorAsTime)
		if err != nil {
			c.String(500, err.Error())
			return err
		}

		threads.PageName = "main"
		return c.Render(200, "main", threads)
	}

	siteContext, err := h.threadService.GetThreadsWithCursor(limit, false, &cursorAsTime)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	siteContext.PageName = "main"
	return c.Render(200, "main", siteContext)
}

func (h *TemplateHandler) Ping(c echo.Context) error {
	return c.Render(200, "ping", nil)
}

func (h *TemplateHandler) ListPostsForThread(c echo.Context) error {
	threadID := c.Param("id")

	idAsInt, err := strconv.Atoi(threadID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	posts, err := h.threadService.GetThreadByID(100, 0, idAsInt)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.Render(200, "posts", posts)
}

func (h *TemplateHandler) Post(c echo.Context) error {
	postID := c.Param("id")
	idAsInt, err := strconv.Atoi(postID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	post, err := h.threadService.GetPostByID(idAsInt)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.Render(200, "post", post)
}

func (h *TemplateHandler) ThreadReply(c echo.Context) error {
	values, err := c.FormParams()
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	threadID := values.Get("thread_id")
	threadIDAsInt, err := strconv.Atoi(threadID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	author := values.Get("member_name")
	body := values.Get("body")

	ip := c.RealIP()

	postID, err := h.threadService.NewPost(body, ip, author, threadIDAsInt)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.JSON(200, NewPostResponse{
		PostID: postID,
	})
}

func (h *TemplateHandler) CreateThread(c echo.Context) error {
	body := c.FormValue("body")
	subject := c.FormValue("subject")

	ip := c.RealIP()

	author := c.FormValue("member")

	threadID, err := h.threadService.NewThread(author, ip, body, subject)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.JSON(200, NewThreadResponse{
		ThreadID: threadID,
	})
}

func (h *TemplateHandler) NewThread(c echo.Context) error {
	return c.Render(200, "newthread", nil)
}

func (h *TemplateHandler) PreviewPost(c echo.Context) error {
	position := c.Param("position")

	positionAsInt, err := strconv.Atoi(position)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	values, err := c.FormParams()
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	body := values.Get("body")
	threadID := values.Get("thread_id")
	author := values.Get("member_name")

	threadIDAsInt, err := strconv.Atoi(threadID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	now := time.Now()

	post := domain.Post{
		Text:           body,
		MemberName:     author,
		ThreadID:       threadIDAsInt,
		Timestamp:      &now,
		ThreadPosition: positionAsInt + 1,
	}

	return c.Render(200, "post", post)
}

type GenericResponse struct {
	Message string `json:"message"`
}
