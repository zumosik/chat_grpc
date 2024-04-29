package models

import "github.com/zumosik/grpc_chat_protos/go/chat"

type Msg struct {
	ID     string
	ChatID string
	UserID string
	Text   string

	// TODO: Payload
}

func (m *Msg) ToProto() *chat.Message {
	return &chat.Message{
		ChatID:  m.ChatID,
		Text:    m.Text,
		Payload: nil, // TODO: Add payload
	}
}
