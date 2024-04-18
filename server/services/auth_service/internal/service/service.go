package service

import (
	"auth_service/internal/models"
	"context"
)

type EmailTokenStorage interface {
	CreateEmailToken(ctx context.Context, token, userID string) error
	GetUserIDByToken(ctx context.Context, token string) (string, error)
	DeleteEmailToken(ctx context.Context, token string) error
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
