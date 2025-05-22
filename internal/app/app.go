package app

import (
	"context"
	"log/slog"

	authclient "github.com/AlexMickh/speak-gateway/internal/clients/auth"
	userClient "github.com/AlexMickh/speak-gateway/internal/clients/user"
	"github.com/AlexMickh/speak-gateway/internal/config"
	"github.com/AlexMickh/speak-gateway/internal/server"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
)

type App struct {
	cfg        *config.Config
	server     *server.Server
	authClient *authclient.AuthClient
	userClient *userClient.UserClient
}

func Register(ctx context.Context, cfg *config.Config) *App {
	const op = "app.Register"

	ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

	sl.GetFromCtx(ctx).Info(ctx, "initing auth client")
	authClient, err := authclient.New(cfg.Clients.AuthAddr)
	if err != nil {
		sl.GetFromCtx(ctx).Fatal(ctx, "failed to init auth client", sl.Err(err))
	}

	sl.GetFromCtx(ctx).Info(ctx, "initing user client")
	userClient, err := userClient.New(cfg.Clients.UserAddr)
	if err != nil {
		sl.GetFromCtx(ctx).Fatal(ctx, "failed to init user client", sl.Err(err))
	}

	sl.GetFromCtx(ctx).Info(ctx, "initing server")
	server := server.New(ctx, cfg.Server, authClient, userClient)

	return &App{
		cfg:        cfg,
		server:     server,
		authClient: authClient,
		userClient: userClient,
	}
}

func (a *App) Run(ctx context.Context) {
	go func() {
		const op = "app.Run"

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

		sl.GetFromCtx(ctx).Info(ctx, "starting server")
		err := a.server.Run()
		if err != nil {
			sl.GetFromCtx(ctx).Fatal(ctx, "failed to start server")
		}
	}()
}

func (a *App) GracefulStop(ctx context.Context) {
	const op = "app.GracefulStop"

	ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

	sl.GetFromCtx(ctx).Info(ctx, "stopping server")
	err := a.server.GracefulStop(ctx)
	if err != nil {
		sl.GetFromCtx(ctx).Fatal(ctx, "failed to stop server", sl.Err(err))
	}

	sl.GetFromCtx(ctx).Info(ctx, "stopping auth client")
	a.authClient.Close()

	sl.GetFromCtx(ctx).Info(ctx, "stopping user client")
	a.userClient.Close()
}
