package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/stdlib"
	"goBoard/helpers/auth/keyset"
	"goBoard/helpers/auth/token"
	"goBoard/internal/core/service/authenticationsvc"
	"goBoard/internal/core/service/imagesvc"
	"goBoard/internal/core/service/membersvc"
	"goBoard/internal/core/service/messagesvc"
	"goBoard/internal/core/service/themesvc"
	"goBoard/internal/core/service/threadsvc"
	"goBoard/internal/repos/authenticationrepo"
	"goBoard/internal/repos/imagerepo"
	"goBoard/internal/repos/memberrepo"
	"goBoard/internal/repos/messagerepo"
	"goBoard/internal/repos/themerepo"
	"goBoard/internal/repos/threadrepo"
	"goBoard/internal/transport/handlers/authentication"
	"goBoard/internal/transport/handlers/members"
	"goBoard/internal/transport/handlers/messages"
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

	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"

	"github.com/gorilla/sessions"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	dbURI := getenv("DB_URI")
	if dbURI == "" {
		return fmt.Errorf("DB_URI is required")
	}

	cognitoClientID := getenv("COGNITO_CLIENT_ID")
	if cognitoClientID == "" {
		return fmt.Errorf("COGNITO_CLIENT_ID is required")
	}

	jwksURI := getenv("JWKS_URI")
	if jwksURI == "" {
		return fmt.Errorf("JWKS_URI is required")
	}

	sessionKey := getenv("SESSION_KEY")
	if sessionKey == "" {
		return fmt.Errorf("SESSION_KEY is required")
	}

	pool, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)

	fs := os.DirFS("db/migrations")
	d, err := iofs.New(fs, ".")
	if err != nil {
		return err
	}

	driver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("iofs", d, "railway", driver)
	if err != nil {
		return err
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
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

	ksetCache := keyset.NewKeySetWithCache(jwksURI, 15)
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
	authRepo := authenticationrepo.NewAuthenticationRepo(cognitoClient, cognitoClientID)
	messageRepo := messagerepo.NewMessageRepo(pool)
	themeRepo := themerepo.NewThemeRepo(pool)

	threadService := threadsvc.NewThreadService(threadRepo, memberRepo, sugar, 50)
	memberService := membersvc.NewMemberService(memberRepo, sugar)
	authService := authenticationsvc.NewAuthenticationService(authRepo, memberRepo, sugar)
	imageService := imagesvc.NewImageService(imageRepo, *sugar)
	messageService := messagesvc.NewMessageService(messageRepo, memberRepo, sugar, 20)
	themeService := themesvc.NewThemeService(themeRepo, sugar)

	threadsHandler := threads.NewHandler(threadService, memberService, imageService, themeService, jwtMiddleware, sugar)
	authHandler := authentication.NewHandler(authService)
	membersHandler := members.NewHandler(threadService, memberService, jwtMiddleware, sugar)
	messagesHandler := messages.NewHandler(messageService, memberService, imageService, jwtMiddleware, sugar)

	r := chi.NewRouter()

	store := sessions.NewCookieStore([]byte(sessionKey))
	r.Use(session.SessionMiddleware(store))

	//r.Use(middleware.RealIP)

	threadsHandler.Register(r)
	authHandler.Register(r)
	membersHandler.Register(r)
	messagesHandler.Register(r)

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("cmd/public"))))

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

	log.Println("** starting bco on port 80 **")
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-serverCtx.Done()

	return nil
}
