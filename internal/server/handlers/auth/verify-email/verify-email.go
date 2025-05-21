package verifyemail

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/AlexMickh/speak-gateway/pkg/utils/render"
)

type EmailVerifier interface {
	VerifyEmail(ctx context.Context, id string) error
}

type Request struct {
	ID string `json:"id"`
}

func New(auth EmailVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.auth.verify-email.New"

		ctx := r.Context()

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "not valid request", sl.Err(err))
			render.JSON(w, http.StatusBadRequest, response.Error("not valid request"))
			return
		}

		err = auth.VerifyEmail(ctx, req.ID)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to verify email", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, response.Error("failed to verify email"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
