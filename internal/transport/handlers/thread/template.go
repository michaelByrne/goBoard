package thread

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
)

type TemplateHandler struct {
	threadService      ports.ThreadService
	memberService      ports.MemberService
	defaultThreadLimit int
}

func NewTemplateHandler(threadService ports.ThreadService, memberService ports.MemberService, defaultThreadLimit int) *TemplateHandler {
	return &TemplateHandler{
		threadService:      threadService,
		memberService:      memberService,
		defaultThreadLimit: defaultThreadLimit,
	}
}

func (h *TemplateHandler) Register(e *echo.Echo) {
	//e.GET("/", h.ListFirstPageThreads)
	e.GET("/", h.ListThreads)
	e.GET("/thread/list", h.ListThreads)
	e.GET("/thread/view/:id", h.ListPostsForThread)
	e.POST("/thread/list/nav", h.ThreadListNav)
	e.GET("/thread/post/:id/:position", h.Post)
	e.POST("/threads", h.Threads)
	e.POST("/thread", h.Thread)
	e.GET("/ping", h.Ping)
	//e.POST("/thread/reply", h.ThreadReply)
	//e.POST("/thread/create", h.CreateThread)
	e.GET("/thread/create", h.NewThread)
	e.POST("/thread/previewpost/:position", h.PreviewPost)
}

func (h *TemplateHandler) ListThreads(c echo.Context) error {
	cursor := c.QueryParams().Get("cursor")
	var cursorAsTime time.Time
	var err error
	if cursor == "" {
		cursorAsTime = time.Date(9999, 1, 1, 1, 1, 1, 1, time.UTC)
	} else {
		cursorAsTime, err = time.Parse(time.RFC3339, cursor)
		if err != nil {
			c.String(500, err.Error())
			return err
		}
	}

	reverse := c.QueryParams().Get("reverse")
	var reverseAsBool bool
	if reverse == "" {
		reverseAsBool = false
	} else {
		reverseAsBool, err = strconv.ParseBool(reverse)
		if err != nil {
			c.String(500, err.Error())
			return err
		}
	}

	args := ThreadsPassThroughArgs{
		Cursor:  &cursorAsTime,
		Reverse: reverseAsBool,
	}

	return c.Render(200, "main", args)
}

//func (h *TemplateHandler) ListFirstPageThreads(c echo.Context) error {
//	siteContext, err := h.threadService.GetThreadsWithCursorForward(h.defaultThreadLimit, true, nil)
//	if err != nil {
//		c.String(500, err.Error())
//		return err
//	}
//	siteContext.ThreadPage.PageNum = 0
//	siteContext.PageName = "main"
//	return c.Render(200, "main", siteContext)
//}
//
//func (h *TemplateHandler) ListThreads(c echo.Context) error {
//	cursor := c.QueryParams().Get("cursor")
//	cursorAsTime, err := time.Parse(time.RFC3339, cursor)
//	if err != nil {
//		c.String(500, err.Error())
//		return err
//	}
//
//	reverse := c.QueryParams().Get("reverse")
//	reverseAsBool, err := strconv.ParseBool(reverse)
//	if err != nil {
//		c.String(500, err.Error())
//		return err
//	}
//
//	if reverseAsBool {
//		siteContext, err := h.threadService.GetThreadsWithCursorReverse(h.defaultThreadLimit, &cursorAsTime)
//		if err != nil {
//			c.String(500, err.Error())
//			return err
//		}
//
//		siteContext.PageName = "main"
//		return c.Render(200, "main", siteContext)
//	}
//
//	siteContext, err := h.threadService.GetThreadsWithCursorForward(h.defaultThreadLimit, false, &cursorAsTime)
//	if err != nil {
//		c.String(500, err.Error())
//		return err
//	}
//
//	siteContext.PageName = "main"
//	return c.Render(200, "main", siteContext)
//}

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

	thread, err := h.threadService.GetThreadByID(100, 0, idAsInt)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	for idx, post := range thread.Posts {
		thread.Posts[idx].HtmlBody, err = h.threadService.ConvertPostBodyBbcodeToHtml(post.Body)
		if err != nil {
			c.String(500, err.Error())
			return err
		}
	}

	return c.Render(200, "posts", thread)
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
	htmlBody, err := h.threadService.ConvertPostBodyBbcodeToHtml(post.Body)
	if err != nil {
		c.String(500, err.Error())
		return err
	}
	post.HtmlBody = htmlBody

	return c.Render(200, "post", post)
}

func (h *TemplateHandler) Thread(c echo.Context) error {
	var thread domain.Thread
	err := c.Bind(&thread)
	if err != nil {
		log.Error()
		c.String(500, err.Error())
		return err
	}

	return c.Render(200, "thread", thread)
}

//func (h *TemplateHandler) ThreadReply(c echo.Context) error {
//	values, err := c.FormParams()
//	if err != nil {
//		c.String(500, err.Error())
//		return err
//	}
//
//	threadID := values.Get("thread_id")
//	threadIDAsInt, err := strconv.Atoi(threadID)
//	if err != nil {
//		c.String(500, err.Error())
//		return err
//	}
//
//	author := values.Get("member_name")
//	body := values.Get("body")
//
//	ip := c.RealIP()
//
//	postID, err := h.threadService.NewPost(body, ip, author, threadIDAsInt)
//	if err != nil {
//		c.String(500, err.Error())
//		return err
//	}
//
//	return c.JSON(200, NewPostResponse{
//		PostID: postID,
//	})
//}

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
	htmlBody, err := h.threadService.ConvertPostBodyBbcodeToHtml(body)
	if err != nil {
		c.String(500, err.Error())
		return err
	}
	threadID := values.Get("thread_id")
	author := values.Get("member_name")

	threadIDAsInt, err := strconv.Atoi(threadID)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	now := time.Now()

	post := domain.ThreadPost{
		HtmlBody:   htmlBody,
		Body:       body,
		MemberName: author,
		ParentID:   threadIDAsInt,
		Timestamp:  &now,
		Position:   positionAsInt + 1,
	}

	return c.Render(200, "post", post)
}

func (h *TemplateHandler) Threads(c echo.Context) error {
	var threads []domain.Thread
	err := c.Bind(&threads)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.Render(200, "threads", threads)
}

func (h *TemplateHandler) ThreadListNav(c echo.Context) error {
	var listNavRequest ListNavRequest
	err := c.Bind(&listNavRequest)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	return c.Render(200, "thread-list-nav", listNavRequest)
}

type GenericResponse struct {
	Message string `json:"message"`
}

type ListNavRequest struct {
	Reverse        bool       `json:"reverse"`
	HasNextPage    bool       `json:"hasNextPage"`
	HasPrevPage    bool       `json:"hasPrevPage"`
	PageCursor     *time.Time `json:"pageCursor"`
	PrevPageCursor *time.Time `json:"prevPageCursor"`
}

type ThreadsPassThroughArgs struct {
	Reverse bool
	Cursor  *time.Time
}
