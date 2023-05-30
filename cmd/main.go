package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"goBoard/internal/core/service/membersvc"
	"goBoard/internal/core/service/messagesvc"
	"goBoard/internal/core/service/threadsvc"
	"goBoard/internal/handlers/member"
	"goBoard/internal/handlers/message"
	"goBoard/internal/handlers/thread"
	"goBoard/internal/repos/memberrepo"
	"goBoard/internal/repos/messagerepo"
	"goBoard/internal/repos/threadrepo"
	"html/template"
	"io"
	"log"
	"os"
	"strconv"
)

////go:embed public/views/*.html
//var templateFiles embed.FS

//
////go:embed public/css/*.css
//var cssFiles embed.FS
//
////go:embed public/js/*.js
//var jsFiles embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	maxThreadLimit := os.Getenv("MAX_THREAD_LIMIT")
	if maxThreadLimit == "" {
		maxThreadLimit = "30"
	}

	maxThreadLimitAsInt, err := strconv.Atoi(maxThreadLimit)
	if err != nil {
		log.Fatal(err)
	}

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Sync()

	sugar := l.Sugar()

	dbURI := "postgres://boardking:test@board-postgres:5432/board?sslmode=disable"
	pool, err := pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		log.Fatal(err)
	}

	memberRepo := memberrepo.NewMemberRepo(pool)
	memberService := membersvc.NewMemberService(memberRepo, sugar)

	threadRepo := threadrepo.NewThreadRepo(pool, maxThreadLimitAsInt)
	threadService := threadsvc.NewThreadService(threadRepo, memberRepo, sugar, maxThreadLimitAsInt)

	messageRepo := messagerepo.NewMessageRepo(pool)
	messageService := messagesvc.NewMessageService(messageRepo, memberRepo, sugar, maxThreadLimitAsInt)

	memberTemplateHandler := member.NewTemplateHandler(threadService, memberService)
	threadTemplateHandler := thread.NewTemplateHandler(threadService, memberService, maxThreadLimitAsInt)
	messageTemplateHandler := message.NewTemplateHandler(messageService)

	threadHTTPHandler := thread.NewHandler(threadService)
	memberHTTPHandler := member.NewHandler(memberService)
	messageHTTPHandler := message.NewHandler(memberService, messageService)

	t := &Template{
		templates: template.Must(template.New("t").Funcs(template.FuncMap{
			"add": func(a, b int) int {
				return a + b
			},
			"sub": func(a, b int) int {
				return a - b
			},
		}).ParseGlob("public/views/*.html")),
	}

	e := echo.New()

	e.Renderer = t
	e.Debug = true

	e.Use(middleware.CORS())
	//e.Use(middleware.Logger())

	threadTemplateHandler.Register(e)
	memberTemplateHandler.Register(e)
	messageTemplateHandler.Register(e)
	threadHTTPHandler.Register(e)
	memberHTTPHandler.Register(e)
	messageHTTPHandler.Register(e)

	e.Static("/static", "public")

	e.Logger.Fatal(e.Start(":8080"))
}
