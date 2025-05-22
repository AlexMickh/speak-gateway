package updateuser

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/speak-gateway/internal/domain/models"
	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/AlexMickh/speak-gateway/pkg/utils/render"
)

type UserUpdater interface {
	UpdateUser(
		ctx context.Context,
		accessToken string,
		username string,
		description string,
		profileImage []byte,
	) (models.User, error)
}

type Response struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
}

func New(user UserUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.user.update-user.New"

		ctx := r.Context()

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

		accessToken, ok := ctx.Value(ctx).(string)
		if !ok || accessToken == "" {
			sl.GetFromCtx(ctx).Error(ctx, "failed to get token from context")
			render.JSON(w, http.StatusBadRequest, response.Error("failed to get access token"))
			return
		}

		username := r.PostFormValue("username")
		description := r.PostFormValue("description")

		file, _, err := r.FormFile("profile_image")
		if err != nil {
			userInfo, err := user.UpdateUser(ctx, accessToken, username, description, nil)
			if err != nil {
				sl.GetFromCtx(ctx).Error(ctx, "failed to update user")
				render.JSON(w, http.StatusInternalServerError, response.Error("failed to update user"))
				return
			}

			responseOK(w, userInfo)
			return
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, file)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to copy image", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, response.Error("server error"))
			return
		}

		userInfo, err := user.UpdateUser(ctx, accessToken, username, description, buf.Bytes())
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to update user")
			render.JSON(w, http.StatusInternalServerError, response.Error("failed to update user"))
			return
		}

		responseOK(w, userInfo)
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
