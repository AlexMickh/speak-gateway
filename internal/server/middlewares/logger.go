package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/speak-gateway/pkg/sl"
)

// FIXME
func Logger(ctx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := sl.GetFromCtx(ctx)

			log = log.WithFields(
				slog.String("request_id", r.Header.Get("request_id")),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("agent", r.UserAgent()),
				slog.String("remote_addr", r.RemoteAddr),
			)
			log.Info(ctx, "request details")

			ctx = context.WithValue(r.Context(), sl.Key, log)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
