package service

import (
	"chat_service/internal/client/auth"
	"chat_service/internal/models"
	"context"
	"errors"
	"github.com/zumosik/grpc_chat_protos/go/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log/slog"
)

const (
	TokenMetadataKey = "token"
)

type MessageStorage interface {
	CreateMessage(ctx context.Context, msg *models.Msg) (*models.Msg, error)
	GetMsgByID(ctx context.Context, id string) (*models.Msg, error)
	DeleteMsgByID(ctx context.Context, id string) error
}

type Service struct {
	l *slog.Logger

	authService *auth.Client
	storage     MessageStorage

	activeUsers []*models.User

	chat.UnimplementedChatServiceServer
}

func New(logger *slog.Logger, authService *auth.Client, storage MessageStorage) *Service {
	return &Service{
		l:           logger,
		authService: authService,
		storage:     storage,
	}
}

func (s *Service) Stream(server chat.ChatService_StreamServer) error {
	ctx := server.Context()

	// 1. authenticate user using token and auth service
	user, err := s.getUser(ctx)
	if err != nil {
		return err
	}

	// 1.1 get user rooms (call to room service)
	// TODO

	// 2. save user as active user
	s.activeUsers = append(s.activeUsers, user)

	for {
		select {
		case <-ctx.Done():
			s.l.Debug("Context is done, closing stream")
			return nil
		default:
			req, err := server.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					s.l.Debug("Client closed the stream")
					return nil
				}
				s.l.Error("Failed to receive a message from stream",
					slog.String("error", err.Error()))
				return status.Error(codes.Internal, "internal error")
			}

			// 1. validate message
			if req.GetMsg() == nil {
				return status.Error(codes.InvalidArgument, "invalid message")
			}
			if req.GetMsg().GetText() == "" ||
				req.GetMsg().GetText() == " " ||
				req.GetMsg().GetChatID() == "" {
				return status.Error(codes.InvalidArgument, "invalid message")
			}
			// 2. check payload
			// TODO
			// 3. store message in db
			msg := &models.Msg{
				ChatID: req.GetMsg().GetChatID(),
				UserID: user.ID,
				Text:   req.GetMsg().GetText(),
			}
			_, err = s.storage.CreateMessage(ctx, msg)
			if err != nil {
				s.l.Error("Failed to store message in db", slog.String("error", err.Error()))
				return status.Error(codes.Internal, "internal error")
			}
			// 3.1 store payload in db (if needed)
			// TODO
			// 4. send message to all subscribers
		}
	}
}

// getUser extracts user token from context and authenticates user using auth service
// returns user if authenticated, otherwise returns error to be sent to client (status.Error)
func (s *Service) getUser(ctx context.Context) (*models.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	token, ok := md[TokenMetadataKey]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	user, err := s.authService.AuthenticateUser(ctx, []byte(token[0]))
	if err != nil {
		s.l.Error("Failed to authenticate user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return user, nil
}
