package postgres

import (
	"auth_service/internal/models"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	// create id
	user.ID = uuid.New().String()

	query :=
		`
INSERT INTO users (id, username, email, encrypted_password, confirmed_email, confirmed_email)
VALUES (:id, :username, :email, :encrypted_password, :confirmed_email, :created_at)
`
	_, err := s.db.NamedExecContext(ctx, query, user)
	return err
}

func (s *Storage) UpdateUser(ctx context.Context, user *models.User) error {
	query :=
		`
UPDATE users SET username = :username,
 	email = :email,
 	encrypted_password = :encrypted_password,
 	confirmed_email = :confirmed_email,
 	created_at = :created_at 
	WHERE id = :id
	`

	_, err := s.db.NamedExecContext(ctx, query, user)
	return err
}

func (s *Storage) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var user models.User
	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		// if here is no user it isn't error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	var user models.User
	err := s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		// if here is no user it isn't error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	var user models.User
	err := s.db.GetContext(ctx, &user, query, username)
	if err != nil {
		// if here is no user it isn't error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
