package app

import (
	"context"
	"log/slog"

	authclient "github.com/AlexMickh/speak-gateway/internal/clients/auth"
	"github.com/AlexMickh/speak-gateway/internal/config"
	"github.com/AlexMickh/speak-gateway/internal/server"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
)

type App struct {
	cfg    *config.Config
	server *server.Server
}

func Register(ctx context.Context, cfg *config.Config) *App {
	const op = "app.Register"

	ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

	sl.GetFromCtx(ctx).Info(ctx, "initing client")
	client, err := authclient.New(cfg.Clients.AuthAddr)
	if err != nil {
		sl.GetFromCtx(ctx).Fatal(ctx, "failed to init auth client", sl.Err(err))
	}

	sl.GetFromCtx(ctx).Info(ctx, "initing server")
	server := server.New(ctx, cfg.Server, client)

	return &App{
		cfg:    cfg,
		server: server,
	}
}

func (a *App) Run(ctx context.Context) {
	const op = "app.Run"

	ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

	sl.GetFromCtx(ctx).Info(ctx, "starting server")
	err := a.server.Run()
	if err != nil {
		sl.GetFromCtx(ctx).Fatal(ctx, "failed to start server")
	}
}

func (a *App) GracefulStop(ctx context.Context) {
	const op = "app.GracefulStop"

	ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

	sl.GetFromCtx(ctx).Info(ctx, "stopping server")
	err := a.server.GracefulStop(ctx)
	if err != nil {
		sl.GetFromCtx(ctx).Fatal(ctx, "failed to stop server", sl.Err(err))
	}
}
