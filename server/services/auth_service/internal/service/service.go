package service

import (
	"auth_service/internal/client/notifications"
	"auth_service/internal/models"
	"context"
	"github.com/zumosik/grpc_chat_protos/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// TODO: add validation

type EmailTokenStorage interface {
	CreateEmailToken(ctx context.Context, token, userID string) error
	GetUserIDByToken(ctx context.Context, token string) (string, error)
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error

	GetUserByID(ctx context.Context, id string) (*models.User, error)

	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type TokenManager interface {
	CreateToken(u *models.User) ([]byte, error)
	ParseToken(token []byte) (string, error)
}

type Service struct {
	st           UserStorage
	stEmailToken EmailTokenStorage

	tokenManager TokenManager
	l            *slog.Logger

	notificationService *notifications.Client

	auth.UnimplementedAuthServiceServer
}

func New(logger *slog.Logger, storage UserStorage, stEmailToken EmailTokenStorage, tokenManager TokenManager, notificationService *notifications.Client) *Service {
	return &Service{
		st:           storage,
		stEmailToken: stEmailToken,

		l:            logger,
		tokenManager: tokenManager,

		notificationService: notificationService,
	}
}

func Register(server *grpc.Server, service *Service) {
	auth.RegisterAuthServiceServer(server, service)
}

func (s *Service) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	// 1. Find user by username
	u, err := s.st.FindUserByUsername(ctx, request.Username)
	if err != nil {
		s.l.Error("Cant find user by username", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if u == nil {
		// 1.1 If user not found by email
		u, err = s.st.FindUserByEmail(ctx, request.Username)
		if err != nil {
			s.l.Error("Cant find user by email", slog.String("error", err.Error()))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	// 1.2 If user not found by email or username return error
	if u == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	// 2. Compare password
	if !u.ComparePassword(request.Password) {
		return nil, status.Error(codes.Unauthenticated, "wrong password")
	}
	// 3. Generate token
	token, err := s.tokenManager.CreateToken(u)
	if err != nil {
		s.l.Error("Cant create token", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth.LoginResponse{
		Success: true,
		Token:   &auth.Token{Token: token},
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, request *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	// 1. Check if username or email already exists
	u, err := s.st.FindUserByEmail(ctx, request.Email)
	if err != nil {
		s.l.Error("Cant find user by email", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if u != nil {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	}
	u, err = s.st.FindUserByUsername(ctx, request.Username)
	if err != nil {
		s.l.Error("Cant find user by username", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if u != nil {
		return nil, status.Error(codes.AlreadyExists, "user with this username already exists")
	}
	// 2. Create user
	user := models.User{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	}
	// 3. Hash password
	err = user.HashPassword()
	if err != nil {
		s.l.Error("Cant hash password", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}
	// 5. Save user
	err = s.st.CreateUser(ctx, &user)
	if err != nil {
		s.l.Error("Cant create user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	// 6. Create token for email confirm
	token, err := s.tokenManager.CreateToken(u)
	if err != nil {
		s.l.Error("Cant create token for email confirm", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	err = s.stEmailToken.CreateEmailToken(ctx, string(token), user.ID)
	if err != nil {
		s.l.Error("Cant save to storage token for email confirm", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	err = s.notificationService.SendEmailConfirmationEmail(ctx, string(token), user.Email)
	if err != nil {
		s.l.Error("Cant use SendEmailConfirmationEmail (notifications service issue)", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth.CreateUserResponse{
		Success: true,
		User:    user.ToAuthUser(),
	}, nil
}

func (s *Service) UpdateUser(ctx context.Context, request *auth.UpdateUserRequest) (*auth.UpdateUserResponse, error) {
	// 1. Get id from token
	id, err := s.tokenManager.ParseToken(request.Token.Token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	user := models.User{
		ID:       id,
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	}

	// 2. Hash password
	err = user.HashPassword()
	if err != nil {
		s.l.Error("Cant hash password", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	err = s.st.UpdateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &auth.UpdateUserResponse{
		Success: true,
		NewUser: user.ToAuthUser(),
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, request *auth.DeleteUserRequest) (*auth.DeleteUserResponse, error) {
	// 1. Get id from token
	id, err := s.tokenManager.ParseToken(request.Token.Token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}
	// 2. Delete user by id
	err = s.st.DeleteUser(ctx, id)
	if err != nil {
		s.l.Error("Cant delete user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth.DeleteUserResponse{
		Success: true,
	}, nil
}

func (s *Service) GetUserByToken(ctx context.Context, request *auth.GetUserByTokenRequest) (*auth.GetUserResponse, error) {
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
	return &auth.GetUserResponse{
		Success:    true,
		PublicUser: u.ToAuthUser(),
	}, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, request *auth.GetUserByEmailRequest) (*auth.GetUserResponse, error) {
	// 1. Find user by email
	u, err := s.st.FindUserByEmail(ctx, request.Email)
	if err != nil {
		s.l.Error("Cant get user by email", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if u == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &auth.GetUserResponse{
		Success:    true,
		PublicUser: u.ToAuthUser(),
	}, nil

}

func (s *Service) GetUserByUsername(ctx context.Context, request *auth.GetUserByUsernameRequest) (*auth.GetUserResponse, error) {
	// 1. Find user by username
	u, err := s.st.FindUserByUsername(ctx, request.Username)
	if err != nil {
		s.l.Error("Cant get user by username", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if u == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &auth.GetUserResponse{
		Success:    true,
		PublicUser: u.ToAuthUser(),
	}, nil
}

func (s *Service) VerifyUser(ctx context.Context, req *auth.VerifyUserRequest) (*auth.VerifyUserResponse, error) {
	// 1. Find userID using token
	userID, err := s.stEmailToken.GetUserIDByToken(ctx, req.VerificationCode)
	if err != nil {
		s.l.Error("Cant get userID by email verification code", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	// 2. Find user by userID
	u, err := s.st.GetUserByID(ctx, userID)
	if err != nil {
		s.l.Error("Cant get user by id", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	if u == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// 3. Update user
	u.ConfirmedEmail = true
	err = s.st.UpdateUser(ctx, u)
	if err != nil {
		s.l.Error("Cant update user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth.VerifyUserResponse{
		Success: true,
	}, nil
}
