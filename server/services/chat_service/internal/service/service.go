package service

import (
	"github.com/zumosik/grpc_chat_protos/go/chat"
)

type Service struct {
	chat.UnimplementedChatServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Stream(chat.ChatService_StreamServer) error {

}
