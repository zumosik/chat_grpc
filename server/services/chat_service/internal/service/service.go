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

type RoomsService interface {
	GetUserRooms(ctx context.Context, userID string) ([]*models.Room, error)
}

type AuthService interface {
	AuthenticateUser(ctx context.Context, token []byte) (*models.User, error)
}

type MessageStorage interface {
	CreateMessage(ctx context.Context, msg *models.Msg) (*models.Msg, error)
	GetMsgByID(ctx context.Context, id string) (*models.Msg, error)
	DeleteMsgByID(ctx context.Context, id string) error
}

type Service struct {
	l *slog.Logger

	authService  AuthService
	roomsService RoomsService

	storage MessageStorage

	userServers map[string]chat.ChatService_StreamServer
	activeUsers map[string][]string

	messagesToSend chan *models.Msg

	chat.UnimplementedChatServiceServer
}

func New(logger *slog.Logger, authService *auth.Client, storage MessageStorage) *Service {
	return &Service{
		l:              logger,
		authService:    authService,
		storage:        storage,
		userServers:    make(map[string]chat.ChatService_StreamServer),
		activeUsers:    make(map[string][]string),
		messagesToSend: make(chan *models.Msg, 100),
	}
}

func (s *Service) SendMessagesLoop() {
	for msg := range s.messagesToSend {
		for _, userID := range s.activeUsers[msg.ChatID] {
			server := s.userServers[userID]
			err := server.Send(&chat.StreamResponse{
				Event: &chat.StreamResponse_ClientMessage{
					ClientMessage: msg.ToProto(),
				},
			})
			if err != nil {
				s.l.Error("Failed to send message to user", slog.String("error", err.Error()))
			}
		}
	}
}

func (s *Service) Stream(server chat.ChatService_StreamServer) error {
	ctx := server.Context()

	// 1. authenticate user using token and auth service
	user, err := s.getUser(ctx)
	if err != nil {
		return err
	}

	// 2. save user as active user
	s.userServers[user.ID] = server

	// 3. get user rooms
	rooms, err := s.roomsService.GetUserRooms(ctx, user.ID)
	if err != nil {
		s.l.Error("Failed to get user rooms", slog.String("error", err.Error()))
		return status.Error(codes.Internal, "internal error")
	}

	// 4. save user as active user in each room
	for _, room := range rooms {
		s.activeUsers[room.ID] = append(s.activeUsers[room.ID], user.ID)
	}

	defer func() {
		delete(s.userServers, user.ID)
	}()

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
			s.messagesToSend <- msg
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
