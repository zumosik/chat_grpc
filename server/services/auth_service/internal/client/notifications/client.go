package notifications

import (
	"context"
	"github.com/zumosik/grpc_chat_protos/go/notifications"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type Client struct {
	l      *slog.Logger
	client notifications.NotificationServiceClient
}

func Connect(log *slog.Logger, add string) (*Client, error) {
	conn, err := grpc.Dial(add, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := notifications.NewNotificationServiceClient(conn)

	c := &Client{
		client: client,
		l:      log,
	}

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
