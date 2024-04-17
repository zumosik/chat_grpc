package postgres

import "context"

func (s *Storage) CreateEmailToken(ctx context.Context, token, userID string) error {
	query := `
INSERT INTO email_confirm_tokens (token, user_id)
VALUES ($1, $2)`
	_, err := s.db.ExecContext(ctx, query, token, userID)
	return err
}

func (s *Storage) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	query := `SELECT user_id FROM email_confirm_tokens WHERE token = $1`
	var userID string
	err := s.db.GetContext(ctx, &userID, query, token)
	return userID, err
}
