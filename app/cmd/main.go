package main

import (
	"context"
	"errors"
	"fmt"
	"goBoard/helpers/auth/keyset"
	"goBoard/helpers/auth/token"
	"goBoard/internal/core/service/authenticationsvc"
	"goBoard/internal/core/service/imagesvc"
	"goBoard/internal/core/service/membersvc"
	"goBoard/internal/core/service/threadsvc"
	"goBoard/internal/repos/authenticationrepo"
	"goBoard/internal/repos/imagerepo"
	"goBoard/internal/repos/memberrepo"
	"goBoard/internal/repos/threadrepo"
	"goBoard/internal/transport/handlers/authentication"
	"goBoard/internal/transport/handlers/members"
	"goBoard/internal/transport/handlers/threads"
	"goBoard/internal/transport/middlewares/jwtauth"
	"goBoard/internal/transport/middlewares/session"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const ClientID = "5tbt4vvhsrkqnif6js92lu51a2"

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
	dbURI := "postgresql://gbd:SlipperyBeef@gbd.cjumgo2q6z1h.us-west-2.rds.amazonaws.com:5432/postgres"
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

	defaultConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		return err
	}

	ksetCache := keyset.NewKeySetWithCache("https://cognito-idp.us-west-2.amazonaws.com/us-west-2_3Ju0VjTAr/.well-known/jwks.json", 15)
	kset, err := ksetCache.NewKeySet()
	if err != nil {
		return err
	}

	verifier := token.NewToken(kset)
	jwtMiddleware := jwtauth.Verify(verifier.Verify, jwtauth.TokenFromCookie, jwtauth.TokenFromHeader)

	cognitoClient := cognito.NewFromConfig(defaultConfig)
	s3Client := s3.NewFromConfig(defaultConfig)
	presignClient := s3.NewPresignClient(s3Client)

	imageRepo := imagerepo.NewImageRepo(s3Client, presignClient, *sugar, "dev-bco-images-private")
	threadRepo := threadrepo.NewThreadRepo(pool, 50)
	memberRepo := memberrepo.NewMemberRepo(pool)
	authRepo := authenticationrepo.NewAuthenticationRepo(cognitoClient, ClientID)

	threadService := threadsvc.NewThreadService(threadRepo, memberRepo, sugar, 50)
	memberService := membersvc.NewMemberService(memberRepo, sugar)
	authService := authenticationsvc.NewAuthenticationService(authRepo, memberRepo, sugar)
	imageService := imagesvc.NewImageService(imageRepo, *sugar)

	threadsHandler := threads.NewHandler(threadService, memberService, imageService, jwtMiddleware, sugar)
	authHandler := authentication.NewHandler(authService)
	membersHandler := members.NewHandler(threadService, memberService, jwtMiddleware, sugar)

	r := chi.NewRouter()

	store := sessions.NewCookieStore([]byte("some-secret-key"))
	r.Use(session.SessionMiddleware(store))

	r.Use(middleware.RealIP)

	threadsHandler.Register(r)
	authHandler.Register(r)
	membersHandler.Register(r)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	server := &http.Server{
		Addr:    ":80",
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
