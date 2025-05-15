package register

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/gin-gonic/gin"
)

type Registerer interface {
	Register(
		ctx context.Context,
		username string,
		email string,
		password string,
		description string,
		avatar []byte,
	) (string, error)
}

type Response struct {
	response.Response
	Id string `json:"id,omitempty"`
}

func New(ctx context.Context, auth Registerer) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "server.handlers.auth.register.New"

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

		username := c.PostForm("username")
		if username == "" {
			sl.GetFromCtx(ctx).Error(ctx, "username is empty")
			c.JSON(http.StatusBadRequest, response.Error("username is required"))
			return
		}

		email := c.PostForm("email")
		if email == "" {
			sl.GetFromCtx(ctx).Error(ctx, "email is empty")
			c.JSON(http.StatusBadRequest, response.Error("email is required"))
			return
		}

		password := c.PostForm("password")
		if password == "" {
			sl.GetFromCtx(ctx).Error(ctx, "password is empty")
			c.JSON(http.StatusBadRequest, response.Error("password is required"))
			return
		}

		description := c.PostForm("description")

		file, _, err := c.Request.FormFile("avatar")
		if err == nil {
			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, file)
			if err != nil {
				sl.GetFromCtx(ctx).Error(ctx, "failed to copy file", sl.Err(err))
				c.JSON(http.StatusInternalServerError, response.Error("server error"))
				return
			}

			id, err := auth.Register(ctx, username, email, password, description, buf.Bytes())
			if err != nil {
				sl.GetFromCtx(ctx).Error(ctx, "failed to register user", sl.Err(err))
				c.JSON(http.StatusInternalServerError, response.Error("failed to register user"))
				return
			}

			responseOK(c, id)
			return
		}

		id, err := auth.Register(ctx, username, email, password, description, nil)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to register user", sl.Err(err))
			c.JSON(http.StatusInternalServerError, response.Error("failed to register user"))
			return
		}

		responseOK(c, id)
	}
}

func responseOK(c *gin.Context, id string) {
	c.JSON(http.StatusCreated, Response{
		Response: response.OK(),
		Id:       id,
	})
}
