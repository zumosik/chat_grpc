package notifications

import (
	"auth_service/internal/config"
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"os"

	"github.com/zumosik/grpc_chat_protos/go/notifications"
	"google.golang.org/grpc"
)

type Client struct {
	l      *slog.Logger
	client notifications.NotificationServiceClient
}

func Connect(log *slog.Logger, addr string, certCfg *config.CertsConfig) (*Client, error) {

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

	client := notifications.NewNotificationServiceClient(conn)

	c := &Client{
		client: client,
		l:      log,
	}

	log.Info("Connected to notifications service", slog.String("address", addr))

	return c, nil
}

// SendEmailConfirmationEmail can take few seconds to complete.
func (c *Client) SendEmailConfirmationEmail(ctx context.Context, token, emailTo string) error {
	resp, err := c.client.SendNotification(ctx, &notifications.NotificationRequest{
		Notification: &notifications.NotificationRequest_ConfirmEmail_{
			ConfirmEmail: &notifications.NotificationRequest_ConfirmEmail{
				Email:            emailTo,
				VerificationCode: token,
			}},
	})
	if err != nil {
		return err
	}

	c.l.Debug("Email sent",
		slog.String("method", "SendEmailConfirmationEmail"),
		slog.String("email", emailTo),
		slog.String("resp status", resp.GetStatus()),
	)

	return nil
}
