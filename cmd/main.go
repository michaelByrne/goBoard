package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"goBoard/internal/core/service/threadsvc"
	"goBoard/internal/handlers/thread"
	"goBoard/internal/repos/threadrepo"
	"html/template"
	"io"
	"log"
)

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

	dbURI := "postgres://boardking:test@localhost:5432/board"
	pool, err := pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		log.Fatal(err)
	}

	threadRepo := threadrepo.NewThreadRepo(pool)
	threadService := threadsvc.NewThreadService(threadRepo, sugar)
	//threadHandler := thread.NewHandler(threadService)
	threadTemplateHandler := thread.NewTemplateHandler(threadService)
	//
	//memberRepo := memberrepo.NewMemberRepo(pool)
	//memberService := membersvc.NewMemberService(memberRepo, sugar)
	//memberHandler := member.NewHandler(memberService)

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()

	e.Renderer = t

	e.Use(middleware.CORS())

	//threadHandler.Register(e)
	//memberHandler.Register(e)

	threadTemplateHandler.Register(e)

	e.Static("/static", "public")

	e.Logger.Fatal(e.Start(":8080"))
}
