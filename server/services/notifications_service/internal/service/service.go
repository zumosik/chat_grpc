package service

import (
	"context"
	"fmt"
	"github.com/zumosik/grpc_chat_protos/go/notifications"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/gomail.v2"
	"log/slog"
)

type Service struct {
	from   string
	dialer *gomail.Dialer
	l      *slog.Logger

	notifications.UnimplementedNotificationServiceServer
}

func New(from string, dialer *gomail.Dialer, log *slog.Logger) *Service {
	return &Service{
		from:   from,
		dialer: dialer,
		l:      log,
	}
}

func Register(server *grpc.Server, service *Service) {
	notifications.RegisterNotificationServiceServer(server, service)
}

func (s *Service) SendNotification(
	ctx context.Context,
	req *notifications.NotificationRequest,
) (
	*notifications.NotificationResponse,
	error,
) {
	if req.GetConfirmEmail() != nil {
		confirmEmail := req.GetConfirmEmail()
		email := confirmEmail.GetEmail()
		userID := confirmEmail.GetUserId()

		body := fmt.Sprintf(`
    <h1>Confirmation code</h1>
    <p>Dear User,</p>
    <p>Thank you for signing up. Here is your confirmation code: </p>
    <p><strong>%s</strong></p>
  `, userID)

		m := gomail.NewMessage()
		m.SetHeader("From", s.from)
		m.SetHeader("To", email)
		m.SetHeader("Subject", "Confirmation Email")
		m.SetBody("text/html", body)

		if err := s.dialer.DialAndSend(m); err != nil {
			s.l.Error("cant send", slog.String("error", err.Error()))
			return &notifications.NotificationResponse{Status: "Email sent unsuccessfully"}, status.Error(codes.Internal, "internal error sending email")
		}

		return &notifications.NotificationResponse{Status: "Email sent successfully"}, nil
	}

	return &notifications.NotificationResponse{Status: "Email sent unsuccessfully"}, status.Error(codes.Unimplemented, "not implemented")

}
