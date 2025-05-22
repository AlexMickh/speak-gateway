package getuser

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/speak-gateway/internal/domain/models"
	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/AlexMickh/speak-gateway/pkg/utils/render"
)

type UserGeter interface {
	GetUser(ctx context.Context, email string) (models.User, error)
}

type Request struct {
	Email string `json:"email"`
}

type Response struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
}

func New(user UserGeter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.user.get-user.New"

		ctx := r.Context()

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "not valid request", sl.Err(err))
			render.JSON(w, http.StatusBadRequest, response.Error("not valid request"))
			return
		}

		user, err := user.GetUser(ctx, req.Email)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to get user", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, response.Error("failed to get user"))
			return
		}

		responseOK(w, user)
	}
}

func responseOK(w http.ResponseWriter, user models.User) {
	render.JSON(w, http.StatusOK, Response{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		Description:     user.Description,
		ProfileImageUrl: user.ProfileImageUrl,
	})
}
