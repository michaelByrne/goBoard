package main

import (
	"context"
	"errors"
	"fmt"
	"goBoard/internal/core/service/threadsvc"
	"goBoard/internal/repos/memberrepo"
	"goBoard/internal/repos/threadrepo"
	"goBoard/internal/transport/handlers/threads"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	err := run(context.Background(), os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	dbURI := "postgres://boardking:test@localhost:5432/board?sslmode=disable"
	pool, err := pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		return err
	}

	l, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer l.Sync()

	sugar := l.Sugar()

	threadRepo := threadrepo.NewThreadRepo(pool, 50)
	memberRepo := memberrepo.NewMemberRepo(pool)

	threadService := threadsvc.NewThreadService(threadRepo, memberRepo, sugar, 50)

	threadsHandler := threads.NewHandler(threadService, sugar)

	r := chi.NewRouter()

	threadsHandler.Register(r)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	serverCtx, serverStopCtx := context.WithCancel(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	log.Println("** starting bco on port 8080 **")
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-serverCtx.Done()

	return nil
}
