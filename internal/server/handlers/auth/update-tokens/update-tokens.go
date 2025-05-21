package updatetokens

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/AlexMickh/speak-gateway/pkg/utils/render"
)

type Updater interface {
	UpdateTokens(ctx context.Context, accessToken, refreshToken string) (string, string, error)
}

type Request struct {
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func New(auth Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.auth.update-tokens.New"

		ctx := r.Context()

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "not valid request", sl.Err(err))
			render.JSON(w, http.StatusBadRequest, response.Error("not valid request"))
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			sl.GetFromCtx(ctx).Error(ctx, "auth header is empty")
			render.JSON(w, http.StatusBadRequest, response.Error("auth header required"))
			return
		}

		if strings.Split(token, " ")[0] != "Bearer" {
			sl.GetFromCtx(ctx).Error(ctx, "wrong token type")
			render.JSON(w, http.StatusBadRequest, response.Error("wrong token type, need Bearer"))
			return
		}

		accessToken, refreshToken, err := auth.UpdateTokens(ctx, token, req.RefreshToken)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to update tokens", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, response.Error("failed to update tokens"))
			return
		}

		responseOK(w, accessToken, refreshToken)
	}
}

func responseOK(w http.ResponseWriter, accessToken, refreshToken string) {
	render.JSON(w, http.StatusCreated, Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
