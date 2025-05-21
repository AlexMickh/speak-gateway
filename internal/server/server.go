package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AlexMickh/speak-gateway/internal/config"
	"github.com/AlexMickh/speak-gateway/internal/server/handlers/auth/login"
	"github.com/AlexMickh/speak-gateway/internal/server/handlers/auth/register"
	updatetokens "github.com/AlexMickh/speak-gateway/internal/server/handlers/auth/update-tokens"
	verifyemail "github.com/AlexMickh/speak-gateway/internal/server/handlers/auth/verify-email"
	"github.com/AlexMickh/speak-gateway/internal/server/middlewares"
)

type AuthClient interface {
	Register(
		ctx context.Context,
		username string,
		email string,
		password string,
		description string,
		avatar []byte,
	) (string, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	VerifyToken(ctx context.Context, accessToken string) error
	UpdateTokens(ctx context.Context, accessToken, refreshToken string) (string, string, error)
	VerifyEmail(ctx context.Context, id string) error
}

type Server struct {
	srv *http.Server
}

func New(ctx context.Context, cfg config.Server, authClient AuthClient) *Server {
	mux := http.NewServeMux()

	auth := http.NewServeMux()

	auth.HandleFunc("POST /auth/register", register.New(authClient))
	auth.HandleFunc("POST /auth/login", login.New(authClient))
	auth.HandleFunc("POST /auth/update-tokens", updatetokens.New(authClient))
	auth.HandleFunc("PATCH /auth/verify-email", verifyemail.New(authClient))

	mux.Handle("/auth/", auth)

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: middlewares.Logger(ctx)(mux),
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) Run() error {
	const op = "server.Run"

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) GracefulStop(ctx context.Context) error {
	const op = "server.GracefulStop"

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
