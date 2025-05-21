package register

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/speak-gateway/pkg/api/response"
	"github.com/AlexMickh/speak-gateway/pkg/sl"
	"github.com/AlexMickh/speak-gateway/pkg/utils/render"
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
	Id string `json:"id,omitempty"`
}

func New(auth Registerer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handlers.auth.register.New"

		ctx := r.Context()

		ctx = sl.GetFromCtx(ctx).With(ctx, slog.String("op", op))

		username := r.PostFormValue("username")
		if username == "" {
			sl.GetFromCtx(ctx).Error(ctx, "username is empty")
			render.JSON(w, http.StatusBadRequest, response.Error("username is required"))
			return
		}

		email := r.PostFormValue("email")
		if email == "" {
			sl.GetFromCtx(ctx).Error(ctx, "email is empty")
			render.JSON(w, http.StatusBadRequest, response.Error("email is required"))
			return
		}

		password := r.PostFormValue("password")
		if password == "" {
			sl.GetFromCtx(ctx).Error(ctx, "password is empty")
			render.JSON(w, http.StatusBadRequest, response.Error("password is required"))
			return
		}

		description := r.PostFormValue("description")
		if description == "" {
			sl.GetFromCtx(ctx).Error(ctx, "description is empty")
			render.JSON(w, http.StatusBadRequest, response.Error("description is required"))
			return
		}

		file, _, err := r.FormFile("avatar")
		if err == nil {
			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, file)
			if err != nil {
				sl.GetFromCtx(ctx).Error(ctx, "failed to copy file", sl.Err(err))
				render.JSON(w, http.StatusInternalServerError, response.Error("server error"))
				return
			}

			id, err := auth.Register(ctx, username, email, password, description, buf.Bytes())
			if err != nil {
				sl.GetFromCtx(ctx).Error(ctx, "failed to register user", sl.Err(err))
				render.JSON(w, http.StatusInternalServerError, response.Error("failed to register user"))
				return
			}

			responseOK(w, id)
			return
		}

		id, err := auth.Register(ctx, username, email, password, description, nil)
		if err != nil {
			sl.GetFromCtx(ctx).Error(ctx, "failed to register user", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, response.Error("failed to register user"))
			return
		}

		responseOK(w, id)
	}
}

func responseOK(w http.ResponseWriter, id string) {
	render.JSON(w, http.StatusCreated, Response{
		Id: id,
	})
}
