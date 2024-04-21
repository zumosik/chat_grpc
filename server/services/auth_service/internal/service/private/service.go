package private

import (
	"auth_service/internal/service"
	"context"
	"github.com/zumosik/grpc_chat_protos/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type Service struct {
	st service.UserStorage

	tokenManager service.TokenManager
	l            *slog.Logger

	auth.UnimplementedPrivateServiceServer
}

func New(logger *slog.Logger, storage service.UserStorage, tokenManager service.TokenManager) *Service {
	return &Service{
		st: storage,

		l:            logger,
		tokenManager: tokenManager,
	}
}

func Register(server *grpc.Server, service *Service) {
	auth.RegisterPrivateServiceServer(server, service)
}

func (s *Service) GetUserByToken(ctx context.Context, request *auth.PrivateGetUserByTokenRequest) (*auth.PrivateGetUserResponse, error) {
	// 1. Parse token
	id, err := s.tokenManager.ParseToken(request.Token.Token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}
	// 2. Find user by id
	u, err := s.st.GetUserByID(ctx, id)
	if err != nil {
		s.l.Error("Cant get user by id", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if u == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &auth.PrivateGetUserResponse{
		Success: true,
		User: &auth.User{
			Id:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Verified: u.ConfirmedEmail,
		},
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, req *auth.PrivateGetUserByIDRequest) (*auth.PrivateGetUserResponse, error) {
	u, err := s.st.GetUserByID(ctx, req.GetId())
	if err != nil {
		s.l.Error("Cant get user by id", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if u == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &auth.PrivateGetUserResponse{
		Success: true,
		User: &auth.User{
			Id:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Verified: u.ConfirmedEmail,
		},
	}, nil
}
