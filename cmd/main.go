package main

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4/pgxpool"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"goBoard/helpers/auth"
	"goBoard/internal/core/service/authenticationsvc"
	"goBoard/internal/core/service/membersvc"
	"goBoard/internal/core/service/messagesvc"
	"goBoard/internal/core/service/threadsvc"
	"goBoard/internal/repos/authenticationrepo"
	"goBoard/internal/repos/memberrepo"
	"goBoard/internal/repos/messagerepo"
	"goBoard/internal/repos/threadrepo"
	"goBoard/internal/transport/handlers/authentication"
	"goBoard/internal/transport/handlers/member"
	"goBoard/internal/transport/handlers/message"
	"goBoard/internal/transport/handlers/thread"
	"goBoard/internal/transport/middlewares/session"
	"html/template"
	"io"
	"log"
	"os"
	"strconv"
)

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
	authenticationService := authenticationsvc.NewAuthenticationService(authenticationRepo, memberRepo, sugar)

	memberTemplateHandler := member.NewTemplateHandler(threadService, memberService)
	threadTemplateHandler := thread.NewTemplateHandler(threadService, memberService, maxThreadLimitAsInt)
	messageTemplateHandler := message.NewTemplateHandler(messageService)
	authenticationTemplateHandler := authentication.NewTemplateHandler(authenticationService)

	threadHTTPHandler := thread.NewHandler(threadService, maxThreadLimitAsInt)
	memberHTTPHandler := member.NewHandler(memberService)
	messageHTTPHandler := message.NewHandler(memberService, messageService)
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

	store := sessions.NewCookieStore([]byte("secret"))
	e.Use(session.Middleware(store))

	e.Use(echojwt.WithConfig(echojwt.Config{
		//NewClaimsFunc: auth.GetJWTClaims,
		SigningKey:   []byte(auth.GetJWTSecret()),
		TokenLookup:  "cookie:access-token", // "<source>:<name>"
		ErrorHandler: auth.JWTErrorChecker,
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/login" {
				return true
			}
			return false
		},
	}))

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
