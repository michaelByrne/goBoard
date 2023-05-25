package thread

import (
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"strconv"
	"time"
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
	e.GET("/thread/list", h.ListFirstPageThreads)
	e.GET("/thread/list/:page", h.ListThreads)
	e.GET("/thread/view/:id", h.ListPostsForThread)
	e.GET("/post/:id/:position", h.Post)
	e.GET("/ping", h.Ping)
	e.GET("/thread/create", h.NewThread)
	e.POST("/thread/previewpost/:position", h.PreviewPost)
}

func (h *TemplateHandler) ListFirstPageThreads(c echo.Context) error {
	threadListLength := 100
	siteContext, err := h.threadService.ListThreads(threadListLength, 0)
	if err != nil {
		c.String(500, err.Error())
		return err
	}
	siteContext.ThreadPage.PageNum = 0
	siteContext.PageName = "main"

	err = c.Render(200, "main", siteContext)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return nil
}

func (h *TemplateHandler) ListThreads(c echo.Context) error {
	threadListLength := 100
	pageNum, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	offset := pageNum * threadListLength
	siteContext, err := h.threadService.ListThreads(threadListLength, offset)
	if err != nil {
		c.String(500, err.Error())
		return err
	}
	siteContext.ThreadPage.PageNum = pageNum
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

	err = c.Render(200, "posts", posts)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return nil
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
