package main

import (
	"context"
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
	"log"
)

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
	threadHandler := thread.NewHandler(threadService)

	memberRepo := memberrepo.NewMemberRepo(pool)
	memberService := membersvc.NewMemberService(memberRepo, sugar)
	memberHandler := member.NewHandler(memberService)

	e := echo.New()

	e.Use(middleware.CORS())

	threadHandler.Register(e)
	memberHandler.Register(e)

	e.Logger.Fatal(e.Start(":8080"))
}
