package rooms

import (
	"chat_service/internal/config"
	"chat_service/internal/models"
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/zumosik/grpc_chat_protos/go/rooms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"os"
)

type PrivateClient struct {
	l      *slog.Logger
	client rooms.PrivateRoomsServiceClient
}

func Connect(log *slog.Logger, addr string, certCfg *config.CertsConfig) (*PrivateClient, error) {
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

	client := rooms.NewPrivateRoomsServiceClient(conn)

	c := &PrivateClient{
		client: client,
		l:      log,
	}

	log.Info("Connected to rooms service", slog.String("address", addr))

	return c, nil
}

func (c *PrivateClient) GetUserRooms(ctx context.Context, userID string) ([]*models.Room, error) {
	resp, err := c.client.GetRoomsByUserID(ctx, &rooms.PrivateGetRoomsByUserIDRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	res := make([]*models.Room, 0, len(resp.GetRooms()))

	for i, r := range resp.GetRooms() {
		ids := make([]string, 0, len(r.GetUsers()))
		for j, u := range r.GetUsers() {
			ids[j] = u.GetId()
		}

		res[i] = &models.Room{
			ID:        r.GetId(),
			Name:      r.GetName(),
			UserIDS:   ids,
			CreatedBy: r.GetCreatedBy().GetId(),
		}
	}

	return res, nil
}
