package main

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"goBoard/internal/core/service/authenticationsvc"
	"goBoard/internal/core/service/membersvc"
	"goBoard/internal/core/service/messagesvc"
	"goBoard/internal/core/service/threadsvc"
	"goBoard/internal/repos/authenticationrepo"
	"goBoard/internal/repos/memberrepo"
	"goBoard/internal/repos/messagerepo"
	"goBoard/internal/repos/threadrepo"
	"goBoard/internal/transport/handlers/authentication"
	member2 "goBoard/internal/transport/handlers/member"
	message2 "goBoard/internal/transport/handlers/message"
	thread2 "goBoard/internal/transport/handlers/thread"
	"goBoard/internal/transport/middlewares/session"
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

	authenticationRepo := authenticationrepo.NewAuthenticationRepo(pool)
	authenticationService := authenticationsvc.NewAuthenticationService(authenticationRepo, sugar)

	memberTemplateHandler := member2.NewTemplateHandler(threadService, memberService)
	threadTemplateHandler := thread2.NewTemplateHandler(threadService, memberService, maxThreadLimitAsInt)
	messageTemplateHandler := message2.NewTemplateHandler(messageService)
	authenticationTemplateHandler := authentication.NewTemplateHandler(authenticationService)

	threadHTTPHandler := thread2.NewHandler(threadService, maxThreadLimitAsInt)
	memberHTTPHandler := member2.NewHandler(memberService)
	messageHTTPHandler := message2.NewHandler(memberService, messageService)
	authenticationHTTPHandler := authentication.NewHTTPHandler(authenticationService)

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
	//e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	//	sugar.Info("request body: ", string(reqBody))
	//	sugar.Info("response body: ", string(resBody))
	//}))
	store := sessions.NewCookieStore([]byte("secret"))
	e.Use(session.Middleware(store))

	threadTemplateHandler.Register(e)
	memberTemplateHandler.Register(e)
	messageTemplateHandler.Register(e)
	authenticationTemplateHandler.Register(e)
	threadHTTPHandler.Register(e)
	memberHTTPHandler.Register(e)
	messageHTTPHandler.Register(e)
	authenticationHTTPHandler.Register(e)

	e.Static("/static", "public")

	e.Logger.Fatal(e.Start(":8080"))
}
