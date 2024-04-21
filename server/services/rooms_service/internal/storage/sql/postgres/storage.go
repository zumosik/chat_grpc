package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"rooms_service/internal/models"
)

type AuthServiceClient interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

// Storage works with postgres db and auth service
type Storage struct {
	db         *sql.DB
	authClient AuthServiceClient
}

func New(db *sql.DB, authClient AuthServiceClient) *Storage {
	return &Storage{db: db, authClient: authClient}
}

func (s *Storage) CreateRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	query := `
	INSERT INTO rooms(id, name, user_ids, created_by_id) VALUES ($1, $2, $3, $4)
`

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	room.ID = id.String()

	userIDs := make([]string, 0, len(room.Users))
	for _, u := range room.Users {
		userIDs = append(userIDs, u.ID)
	}

	_, err = s.db.ExecContext(ctx, query, room.ID, room.Name, userIDs, room.CreatedBy.ID)
	if err != nil {
		return nil, err
	}

	return room, nil

}

func (s *Storage) GetRoom(ctx context.Context, id string) (*models.Room, error) {
	query := `
	SELECT id, name, user_ids, created_by_id FROM rooms WHERE id = $1
`

	var room models.Room
	var userIDs []string

	err := s.db.QueryRowContext(ctx, query, id).Scan(&room.ID, &room.Name, &userIDs, &room.CreatedBy.ID)
	if err != nil {
		return nil, err
	}

	// Get users
	room.Users = make([]*models.User, 0, len(userIDs))
	for _, id := range userIDs {
		u, err := s.authClient.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}
		room.Users = append(room.Users, u)
	}

	return &room, nil
}

func (s *Storage) UpdateRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	query := `
	UPDATE rooms SET name = $1, user_ids = $2, created_by_id = $3 WHERE id = $4
`

	userIDs := make([]string, 0, len(room.Users))
	for _, u := range room.Users {
		userIDs = append(userIDs, u.ID)
	}

	_, err := s.db.ExecContext(ctx, query, room.Name, userIDs, room.CreatedBy.ID, room.ID)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *Storage) DeleteRoom(ctx context.Context, id string) error {
	query := `
	DELETE FROM rooms WHERE id = $1
`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetRoomsByUser(ctx context.Context, u *models.User) ([]*models.Room, error) {
	query := `
	SELECT id, name, user_ids, created_by_id FROM rooms
	WHERE $1 = ANY(user_ids)
	`

	rows, err := s.db.QueryContext(ctx, query, u.ID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	rooms := make([]*models.Room, 0)
	for rows.Next() {
		var room models.Room
		var userIDs []string

		err := rows.Scan(&room.ID, &room.Name, &userIDs, &room.CreatedBy.ID)
		if err != nil {
			return nil, err
		}

		// Get users for each room
		room.Users = make([]*models.User, 0, len(userIDs))
		for _, id := range userIDs {
			u, err := s.authClient.GetUserByID(ctx, id)
			if err != nil {
				return nil, err
			}
			room.Users = append(room.Users, u)
		}

		rooms = append(rooms, &room)
	}

	return rooms, nil
}
