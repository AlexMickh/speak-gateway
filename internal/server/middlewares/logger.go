package middlewares

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/gin-gonic/gin"
)

type logger struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (l *logger) Write(b []byte) (int, error) {
	l.body.Write(b)
	return l.ResponseWriter.Write(b)
}

func RequestLoggingMiddleware(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := &logger{
			body:           bytes.Buffer{},
			ResponseWriter: c.Writer,
		}
		c.Writer = logger

		c.Next()

		log := sl.GetFromCtx(ctx)

		log.WithFields(
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		).Info(ctx, "request details")

		// logger.WithFields(log.Fields{
		//     "status":       ctx.Writer.Status(),
		//     "method":       ctx.Request.Method,
		//     "path":         ctx.Request.URL.Path,
		//     "query_params": ctx.Request.URL.Query(),
		//     "req_body":     string(data),
		//     "res_body":     ginBodyLogger.body.String(),
		// }).Info("request details")
	}
}
