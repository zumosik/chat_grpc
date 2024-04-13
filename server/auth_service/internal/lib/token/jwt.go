package token

import (
	"auth_service/internal/models"
	"github.com/golang-jwt/jwt"
	"time"
)

type Manager struct {
	secretKey      []byte
	expirationTime time.Duration
}

func NewTokenService(secretKey string, expirationTime time.Duration) *Manager {
	return &Manager{secretKey: []byte(secretKey), expirationTime: expirationTime}
}

// CreateToken creates a new JWT token for the user
func (ts *Manager) CreateToken(u *models.User) ([]byte, error) {
	expirationTime := time.Now().Add(ts.expirationTime)

	claims := &jwt.StandardClaims{
		Subject:   u.ID,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ts.secretKey)

	if err != nil {
		return nil, err
	}

	return []byte(tokenString), nil
}

// ParseToken parses the JWT token and returns the user ID
func (ts *Manager) ParseToken(tokenStrBytes []byte) (string, error) {
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(string(tokenStrBytes), claims, func(token *jwt.Token) (interface{}, error) {
		return ts.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	return claims.Subject, nil
}
