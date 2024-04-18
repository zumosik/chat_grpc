package auth

import (
	"chat_service/internal/models"
	"context"
	"github.com/zumosik/grpc_chat_protos/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log/slog"
)

// Client is a client for PRIVATE auth service
type Client struct {
	l      *slog.Logger
	client auth.PrivateServiceClient
}

func Connect(log *slog.Logger, addr, crtPath, keyPath string) (*Client, error) {
	creds, err := credentials.NewClientTLSFromFile(crtPath, keyPath)
	if err != nil {
		log.Error("failed to load credentials", slog.String("error", err.Error()))
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))

	if err != nil {
		return nil, err
	}

	client := auth.NewPrivateServiceClient(conn)

	c := &Client{
		client: client,
		l:      log,
	}

	log.Info("Connected to notifications service", slog.String("address", addr))

	return c, nil
}

func (c *Client) AuthenticateUser(ctx context.Context, token []byte) (*models.User, error) {
	resp, err := c.client.GetUserByToken(ctx, &auth.PrivateGetUserByTokenRequest{
		Token: &auth.Token{Token: token},
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:             resp.User.Id,
		Username:       resp.User.Username,
		Email:          resp.User.Email,
		ConfirmedEmail: resp.User.Verified,
	}, nil
}
