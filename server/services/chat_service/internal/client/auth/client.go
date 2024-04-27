package auth

import (
	"chat_service/internal/config"
	"chat_service/internal/models"
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/zumosik/grpc_chat_protos/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"os"
)

// Client is a client for PRIVATE auth service
type Client struct {
	l      *slog.Logger
	client auth.AuthServiceClient
}

func Connect(log *slog.Logger, addr string, certCfg *config.CertsConfig) (*Client, error) {

	//log.Error("failed to load credentials", slog.String("error", err.Error()))

	// Load CA cert
	pemServerCA, err := os.ReadFile(certCfg.CaPath)
	if err != nil {
		log.Error("failed to load CA cert", slog.String("error", err.Error()))
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		log.Error("failed to append CA cert")
		return nil, err
	}

	clientCert, err := tls.LoadX509KeyPair(certCfg.CertPath, certCfg.KeyPath)
	if err != nil {
		log.Error("failed to load client cert", slog.String("error", err.Error()))
		return nil, err
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(cfg)))

	if err != nil {
		return nil, err
	}

	client := auth.NewAuthServiceClient(conn)

	c := &Client{
		client: client,
		l:      log,
	}

	log.Info("Connected to auth service", slog.String("address", addr))

	return c, nil
}

func (c *Client) AuthenticateUser(ctx context.Context, token []byte) (*models.User, error) {
	resp, err := c.client.GetUserByToken(ctx, &auth.GetUserByTokenRequest{
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
