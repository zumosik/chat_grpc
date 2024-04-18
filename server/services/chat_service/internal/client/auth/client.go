package auth

import (
	"chat_service/internal/models"
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/zumosik/grpc_chat_protos/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log/slog"
	"os"
)

// Client is a client for PRIVATE auth service
type Client struct {
	l      *slog.Logger
	client auth.PrivateServiceClient
}

func Connect(log *slog.Logger, addr, certPath string) (*Client, error) {

	//log.Error("failed to load credentials", slog.String("error", err.Error()))

	f, err := os.Open(certPath)
	if err != nil {
		log.Error("failed to open certificate file", slog.String("error", err.Error()), slog.String("path", certPath))
		return nil, err
	}
	pemServerCA, err := io.ReadAll(f)
	if err != nil {
		log.Error("failed to read certificate file", slog.String("error", err.Error()), slog.String("path", certPath))
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		log.Error("failed to append certificate to pool")
		return nil, err
	}

	cfg := &tls.Config{
		RootCAs: certPool,
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(cfg)))

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
