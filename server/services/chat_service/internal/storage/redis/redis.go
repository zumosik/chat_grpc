package redis

import (
	"chat_service/internal/models"
	"context"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

const (
	ChatIDKey = "ChatID"
	UserIDKey = "UserID"
	TextKey   = "Text"
)

type Storage struct {
	client *redis.Client
}

func New(client *redis.Client) *Storage {
	return &Storage{client: client}
}

func (s *Storage) CreateMessage(ctx context.Context, msg *models.Msg) (*models.Msg, error) {
	// 1. create id
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	msg.ID = id.String()
	// 2. save message in redis
	err = s.client.WithContext(ctx).HMSet(msg.ID, map[string]interface{}{
		ChatIDKey: msg.ChatID,
		UserIDKey: msg.UserID,
		TextKey:   msg.Text,
	}).Err()
	if err != nil {
		return nil, err
	}
	// 3. save payload in redis (if needed)
	// TODO

	return msg, nil
}

func (s *Storage) GetMsgByID(ctx context.Context, id string) (*models.Msg, error) {
	msg, err := s.client.WithContext(ctx).HGetAll(id).Result()
	if err != nil {
		return nil, err
	}
	return &models.Msg{
		ID:     id,
		ChatID: msg[ChatIDKey],
		UserID: msg[UserIDKey],
		Text:   msg[TextKey],
	}, nil
}

func (s *Storage) DeleteMsgByID(ctx context.Context, id string) error {
	return s.client.WithContext(ctx).Del(id).Err()
}
