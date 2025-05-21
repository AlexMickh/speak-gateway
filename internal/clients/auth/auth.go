package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexMickh/speak-gateway/pkg/utils/retry"
	"github.com/AlexMickh/speak-protos/pkg/api/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AuthClient struct {
	conn *grpc.ClientConn
	auth auth.AuthClient
}

func New(addr string) (*AuthClient, error) {
	const op = "grpc.clients.auth.New"

	var conn *grpc.ClientConn
	var authClient auth.AuthClient

	retry.WithDelay(5, 500*time.Millisecond, func() error {
		connect, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		auth := auth.NewAuthClient(conn)

		conn = connect
		authClient = auth

		return nil
	})

	return &AuthClient{
		conn: conn,
		auth: authClient,
	}, nil
}

func (a *AuthClient) Register(
	ctx context.Context,
	username string,
	email string,
	password string,
	description string,
	avatar []byte,
) (string, error) {
	const op = "clients.auth.Register"

	res, err := a.auth.Register(ctx, &auth.RegisterRequest{
		Email:        email,
		Username:     username,
		Password:     password,
		Description:  description,
		ProfileImage: avatar,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res.GetId(), nil
}

func (a *AuthClient) Login(ctx context.Context, email, password string) (string, string, error) {
	const op = "clients.auth.Login"

	res, err := a.auth.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return res.GetAccessToken(), res.GetRefreshToken(), nil
}

func (a *AuthClient) VerifyToken(ctx context.Context, accessToken string) error {
	const op = "clients.auth.VerifyToken"

	_, err := a.auth.VerifyToken(ctx, &auth.VerifyTokenRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthClient) UpdateTokens(ctx context.Context, accessToken, refreshToken string) (string, string, error) {
	const op = "clients.auth.UpdateTokens"

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	res, err := a.auth.UpdateTokens(ctx, &auth.UpdateTokensRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return res.GetAccessToken(), res.GetRefreshToken(), nil
}

func (a *AuthClient) VerifyEmail(ctx context.Context, id string) error {
	const op = "clients.auth.VerifyEmail"

	_, err := a.auth.VerifyEmail(ctx, &auth.VerifyEmailRequest{
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
