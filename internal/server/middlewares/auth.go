package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/AlexMickh/speak-gateway/pkg/utils/render"
)

type TokenVerifier interface {
	VerifyToken(ctx context.Context, accessToken string) error
}

func Auth(auth TokenVerifier) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "server.middlewares.auth.New"

			ctx := r.Context()

			ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

			accessToken := r.Header.Get("Authorization")
			if accessToken == "" {
				sl.GetFromCtx(ctx).Error(ctx, "auth header is empty")
				render.JSON(w, http.StatusBadRequest, response.Error("auth header required"))
				return
			}

			if strings.Split(accessToken, " ")[0] != "Bearer" {
				sl.GetFromCtx(ctx).Error(ctx, "wrong token type")
				render.JSON(w, http.StatusBadRequest, response.Error("wrong token type, need Bearer"))
			}

			err := auth.VerifyToken(ctx, accessToken)
			if err != nil {
				sl.GetFromCtx(ctx).Error(ctx, "can't verify access token", sl.Err(err))
				render.JSON(w, http.StatusUnauthorized, response.Error("user unauthorized"))
				return
			}

			ctx = context.WithValue(r.Context(), "auth", accessToken)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
