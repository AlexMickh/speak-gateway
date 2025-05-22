package user

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexMickh/speak-gateway/internal/domain/models"
	"github.com/AlexMickh/speak-gateway/pkg/utils/retry"
	"github.com/AlexMickh/speak-protos/pkg/api/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type UserClient struct {
	conn *grpc.ClientConn
	user user.UserClient
}

func New(addr string) (*UserClient, error) {
	const op = "clients.user.New"

	var conn *grpc.ClientConn
	var userClient user.UserClient

	err := retry.WithDelay(5, 500*time.Millisecond, func() error {
		connect, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		user := user.NewUserClient(conn)

		conn = connect
		userClient = user

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &UserClient{
		conn: conn,
		user: userClient,
	}, nil
}

func (u *UserClient) GetUser(ctx context.Context, email string) (models.User, error) {
	const op = "clients.user.GetUser"

	res, err := u.user.GetUser(ctx, &user.GetUserRequest{
		Email: email,
	})
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.User{
		ID:              res.GetId(),
		Email:           res.GetEmail(),
		Username:        res.GetUsername(),
		Description:     res.GetDescription(),
		ProfileImageUrl: res.GetProfileImageUrl(),
	}, nil
}

func (u *UserClient) UpdateUser(
	ctx context.Context,
	accessToken string,
	username string,
	description string,
	profileImage []byte,
) (models.User, error) {
	const op = "clients.user.UpdateUser"

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	res, err := u.user.UpdateUserInfo(ctx, &user.UpdateUserInfoRequest{
		Username:     username,
		Description:  description,
		ProfileImage: profileImage,
	})
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.User{
		ID:              res.GetId(),
		Email:           res.GetEmail(),
		Username:        res.GetUsername(),
		Description:     res.GetDescription(),
		ProfileImageUrl: res.GetProfileImageUrl(),
	}, nil
}

func (u *UserClient) DeleteUser(ctx context.Context, id string) error {
	const op = "clients.user.DeleteUser"

	_, err := u.user.DeleteUser(ctx, &user.DeleteUserRequest{
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserClient) Close() {
	u.conn.Close()
}
