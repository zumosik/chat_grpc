package service

import (
	"chat_service/internal/client/auth"
	"github.com/zumosik/grpc_chat_protos/go/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type Service struct {
	l *slog.Logger

	authService *auth.Client

	chat.UnimplementedChatServiceServer
}

func New(logger *slog.Logger, authService *auth.Client) *Service {
	return &Service{l: logger, authService: authService}
}

func (s *Service) Stream(server chat.ChatService_StreamServer) error {
	ctx := server.Context()
	for {
		select {
		case <-ctx.Done():
			s.l.Debug("Context is done, closing stream")
			return nil
		default:
			req, err := server.Recv()
			if err != nil {
				s.l.Error("Failed to receive a message from stream",
					slog.String("error", err.Error()))
				return status.Error(codes.Internal, "internal error")
			}

			// 1. authenticate user using token and auth service
			user, err := s.authService.AuthenticateUser(ctx, req.GetToken().GetToken())
			if err != nil {
				s.l.Error("Failed to authenticate user", slog.String("error", err.Error()))
				return status.Error(codes.Unauthenticated, "unauthenticated")
			}
			if user == nil {
				return status.Error(codes.Unauthenticated, "unauthenticated")
			}
			// 2. validate message
			if req.GetMsg() == nil {
				return status.Error(codes.InvalidArgument, "invalid message")
			}
			if req.GetMsg().GetText() == "" ||
				req.GetMsg().GetText() == " " ||
				req.GetMsg().GetChatID() == "" {
				return status.Error(codes.InvalidArgument, "invalid message")
			}
			// 3. check payload
			// TODO
			// 4. store message in db
			// 4.1 store payload in db (if needed)
			// TODO
			// 5. send message to all subscribers
		}
	}
}
