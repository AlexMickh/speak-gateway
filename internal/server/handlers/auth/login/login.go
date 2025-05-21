package login

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/AlexMickh/speak-gateway/pkg/utils/render"
)

type Loginer interface {
	Login(ctx context.Context, email, password string) (string, string, error)
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	response.Response
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func New(auth Loginer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.auth.login.New"

		ctx := r.Context()

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "not valid request", sl.Err(err))
			render.JSON(w, http.StatusBadRequest, response.Error("not valid request"))
			return
		}

		accessToken, refreshToken, err := auth.Login(ctx, req.Email, req.Password)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to login user", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, response.Error("failed to login user"))
			return
		}

		responseOK(w, accessToken, refreshToken)
	}
}

func responseOK(w http.ResponseWriter, accessToken, refreshToken string) {
	render.JSON(w, http.StatusOK, Response{
		Response:     response.OK(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
