package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/zumosik/grpc_chat_protos/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"os"
	"rooms_service/internal/config"
	"rooms_service/internal/models"
)

// Client is a client for PRIVATE auth service
type Client struct {
	l      *slog.Logger
	client auth.PrivateServiceClient
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

	client := auth.NewPrivateServiceClient(conn)

	c := &Client{
		client: client,
		l:      log,
	}

	log.Info("Connected to private auth service", slog.String("address", addr))

	return c, nil
}

func (c *Client) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
	resp, err := c.client.GetUserByToken(ctx, &auth.PrivateGetUserByTokenRequest{
		Token: &auth.Token{Token: []byte(token)},
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       resp.User.Id,
		Username: resp.User.Username,
		Email:    resp.User.Email,
		Verified: resp.User.Verified,
	}, nil
}

func (c *Client) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	resp, err := c.client.GetUserByID(ctx, &auth.PrivateGetUserByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       resp.User.Id,
		Username: resp.User.Username,
		Email:    resp.User.Email,
		Verified: resp.User.Verified,
	}, nil
}