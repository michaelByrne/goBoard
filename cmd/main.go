package main

import (
	"context"
	"embed"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"goBoard/internal/core/service/membersvc"
	"goBoard/internal/core/service/threadsvc"
	"goBoard/internal/handlers/member"
	"goBoard/internal/handlers/thread"
	"goBoard/internal/repos/memberrepo"
	"goBoard/internal/repos/threadrepo"
	"html/template"
	"io"
	"log"
)

//go:embed public/views/*.html
var templateFiles embed.FS

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
	threadRepo := threadrepo.NewThreadRepo(pool)
	threadService := threadsvc.NewThreadService(threadRepo, memberRepo, sugar)

	memberTemplateHandler := member.NewTemplateHandler(threadService, memberService)
	threadTemplateHandler := thread.NewTemplateHandler(threadService, memberService)

	threadHTTPHandler := thread.NewHandler(threadService)

	t := &Template{
		templates: template.Must(template.New("t").Funcs(template.FuncMap{
			"add": func(a, b int) int {
				return a + b
			},
			"sub": func(a, b int) int {
				return a - b
			},
		}).ParseFS(templateFiles, "public/views/*.html")),
	}

	e := echo.New()

	e.Renderer = t
	e.Debug = true

	e.Use(middleware.CORS())

	threadTemplateHandler.Register(e)
	memberTemplateHandler.Register(e)
	threadHTTPHandler.Register(e)

	e.Static("/static", "public")

	e.Logger.Fatal(e.Start(":8080"))
}
