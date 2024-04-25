package service

import (
	"context"
	"github.com/zumosik/grpc_chat_protos/go/rooms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"rooms_service/internal/models"
)

type PrivateService struct {
	l       *slog.Logger
	storage RoomStorage

	rooms.UnimplementedPrivateRoomsServiceServer
}

func RegisterPrivate(server *grpc.Server, service *PrivateService) {
	rooms.RegisterPrivateRoomsServiceServer(server, service)
}

func (p *PrivateService) GetRoomsByUserID(ctx context.Context, request *rooms.PrivateGetRoomsByUserIDRequest) (*rooms.PrivateGetRoomsByUserResponse, error) {
	u := models.User{ // dont need other fields
		ID: request.GetUserId(),
	}

	userRooms, err := p.storage.GetRoomsByUser(ctx, &u)
	if err != nil {
		p.l.Error("Cant get rooms", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "Cant get rooms")
	}

	roomsRes := make([]*rooms.PrivateRoom, 0, len(userRooms))
	for i, r := range userRooms {
		users := make([]*rooms.User, 0, len(r.Users))
		for j, u := range r.Users {
			users[j] = &rooms.User{
				Id:       u.ID,
				Username: u.Username,
				Email:    u.Email,
				Verified: u.Verified,
			}
		}

		roomsRes[i] = &rooms.PrivateRoom{
			Id:    r.ID,
			Name:  r.Name,
			Users: users,
			CreatedBy: &rooms.User{
				Id:       r.CreatedBy.ID,
				Username: r.CreatedBy.Username,
				Email:    r.CreatedBy.Email,
				Verified: r.CreatedBy.Verified,
			},
		}
	}

	return &rooms.PrivateGetRoomsByUserResponse{
		Success: true,
		Rooms:   roomsRes,
	}, nil
}

func NewPrivate(l *slog.Logger, storage RoomStorage) *PrivateService {
	return &PrivateService{l: l, storage: storage}
}
