package login

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/gin-gonic/gin"
)

type Loginer interface {
	Login(ctx context.Context, email, password string) (string, string, error)
}

type Request struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Response struct {
	response.Response
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func New(ctx context.Context, auth Loginer) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "server.handlers.auth.login.New"

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

		var req Request
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "not valid request", sl.Err(err))
			c.JSON(http.StatusBadGateway, response.Error("not valid request"))
			return
		}

		accessToken, refreshToken, err := auth.Login(ctx, req.Email, req.Password)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to login user", sl.Err(err))
			c.JSON(http.StatusInternalServerError, response.Error("failed to login user"))
			return
		}

		responseOK(c, accessToken, refreshToken)
	}
}

func responseOK(c *gin.Context, accessToken, refreshToken string) {
	c.JSON(http.StatusOK, Response{
		Response:     response.OK(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
