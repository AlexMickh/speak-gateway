package deleteuser

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/AlexMickh/speak-gateway/pkg/utils/render"
)

type UserDeleter interface {
	DeleteUser(ctx context.Context, id string) error
}

type Request struct {
	ID string `json:"id"`
}

func New(user UserDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.user.delete-user.New"

		ctx := r.Context()

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "not valid request", sl.Err(err))
			render.JSON(w, http.StatusBadRequest, response.Error("not valid request"))
			return
		}

		err = user.DeleteUser(ctx, req.ID)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to delete user", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, response.Error("failed to delete user"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
