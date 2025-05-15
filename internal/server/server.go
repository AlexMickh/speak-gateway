package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AlexMickh/speak-gateway/internal/config"
	"github.com/AlexMickh/speak-gateway/internal/server/handlers/auth/login"
	"github.com/AlexMickh/speak-gateway/internal/server/handlers/auth/register"
	"github.com/AlexMickh/speak-gateway/internal/server/middlewares"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/gin-gonic/gin"
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
}

type Server struct {
	srv *http.Server
}

func New(ctx context.Context, cfg config.Server, authClient AuthClient) *Server {
	r := gin.Default()

	r.Use(middlewares.RequestLoggingMiddleware(ctx))
	r.Use(gin.LoggerWithWriter(sl.GetFromCtx(ctx).Writer()))

	r.POST("/auth/register", register.New(ctx, authClient))
	r.POST("/auth/login", login.New(ctx, authClient))

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r.Handler(),
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
