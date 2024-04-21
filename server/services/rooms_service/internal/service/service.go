package service

import (
	"context"
	"github.com/zumosik/grpc_chat_protos/go/rooms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"rooms_service/internal/interceptor"
	"rooms_service/internal/models"
)

// TODO: RemoveFromRoom

type RoomStorage interface {
	CreateRoom(ctx context.Context, room *models.Room) (*models.Room, error)
	GetRoom(ctx context.Context, id string) (*models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) (*models.Room, error)
	DeleteRoom(ctx context.Context, id string) error

	GetRoomsByUser(ctx context.Context, u *models.User) ([]*models.Room, error)
}

type Service struct {
	l *slog.Logger

	storage RoomStorage

	rooms.UnimplementedRoomServiceServer
}

func New(l *slog.Logger, storage RoomStorage) *Service {
	return &Service{l: l, storage: storage}
}

func Register(server *grpc.Server, service *Service) {
	rooms.RegisterRoomServiceServer(server, service)
}

func (s *Service) CreateRoom(ctx context.Context, req *rooms.CreateRoomRequest) (*rooms.CreateRoomResponse, error) {
	u, ok := ctx.Value(interceptor.UserContextKey).(*models.User)
	if !ok {
		s.l.Error("Cant get user from context")
		return nil, status.Error(codes.Internal, "Cant authenticate user")
	}

	room := &models.Room{
		Name:      req.Name,
		Users:     []*models.User{u},
		CreatedBy: u,
	}

	roomResp, err := s.storage.CreateRoom(ctx, room)
	if err != nil {
		s.l.Error("Cant create room", slog.String("error", err.Error()))
	}

	return &rooms.CreateRoomResponse{Room: roomResp.ToProto()}, nil
}

func (s *Service) GetRoom(ctx context.Context, req *rooms.GetRoomRequest) (*rooms.GetRoomResponse, error) {
	u, ok := ctx.Value(interceptor.UserContextKey).(*models.User)
	if !ok {
		s.l.Error("Cant get user from context")
		return nil, status.Error(codes.Internal, "Cant authenticate user")
	}

	room, err := s.storage.GetRoom(ctx, req.RoomId)
	if err != nil {
		s.l.Error("Cant get room", slog.String("error", err.Error()))
		return nil, status.Error(codes.NotFound, "Room not found")
	}

	if !room.HasUser(u) {
		return nil, status.Error(codes.PermissionDenied, "User is not in the room")
	}

	return &rooms.GetRoomResponse{Room: room.ToProto()}, nil
}

func (s *Service) UpdateRoom(ctx context.Context, req *rooms.UpdateRoomRequest) (*rooms.UpdateRoomResponse, error) {
	u, ok := ctx.Value(interceptor.UserContextKey).(*models.User)
	if !ok {
		s.l.Error("Cant get user from context")
		return nil, status.Error(codes.Internal, "Cant authenticate user")
	}

	room, err := s.storage.GetRoom(ctx, req.RoomId)
	if err != nil {
		s.l.Error("Cant get room", slog.String("error", err.Error()))
		return nil, status.Error(codes.NotFound, "Room not found")
	}

	if room.CreatedBy.ID != u.ID {
		return nil, status.Error(codes.PermissionDenied, "User is not the creator of the room")
	}

	// add fields here to update
	// users cant be updated here, use AddToRoom and RemoveFromRoom
	room.Name = req.Name

	roomResp, err := s.storage.UpdateRoom(ctx, room)
	if err != nil {
		s.l.Error("Cant update room", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "Cant update room")
	}

	return &rooms.UpdateRoomResponse{Room: roomResp.ToProto()}, nil
}

func (s *Service) DeleteRoom(ctx context.Context, req *rooms.DeleteRoomRequest) (*rooms.DeleteRoomResponse, error) {
	u, ok := ctx.Value(interceptor.UserContextKey).(*models.User)
	if !ok {
		s.l.Error("Cant get user from context")
		return nil, status.Error(codes.Internal, "Cant authenticate user")
	}

	room, err := s.storage.GetRoom(ctx, req.RoomId)
	if err != nil {
		s.l.Error("Cant get room", slog.String("error", err.Error()))
		return nil, status.Error(codes.NotFound, "Room not found")
	}

	if room.CreatedBy.ID != u.ID {
		return nil, status.Error(codes.PermissionDenied, "User is not the creator of the room")
	}

	err = s.storage.DeleteRoom(ctx, req.RoomId)
	if err != nil {
		s.l.Error("Cant delete room", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "Cant delete room")
	}

	return &rooms.DeleteRoomResponse{}, nil
}

func (s *Service) AddToRoom(ctx context.Context, req *rooms.AddToRoomRequest) (*rooms.AddToRoomResponse, error) {
	u, ok := ctx.Value(interceptor.UserContextKey).(*models.User)
	if !ok {
		s.l.Error("Cant get user from context")
		return nil, status.Error(codes.Internal, "Cant authenticate user")
	}

	room, err := s.storage.GetRoom(ctx, req.RoomId)
	if err != nil {
		s.l.Error("Cant get room", slog.String("error", err.Error()))
		return nil, status.Error(codes.NotFound, "Room not found")
	}

	if room.HasUser(u) {
		return nil, status.Error(codes.AlreadyExists, "User is already in the room")
	}

	room.Users = append(room.Users, u)

	roomResp, err := s.storage.UpdateRoom(ctx, room)
	if err != nil {
		s.l.Error("Cant update room", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "Cant update room")
	}

	return &rooms.AddToRoomResponse{Room: roomResp.ToProto()}, nil
}

func (s *Service) DeleteFromRoom(ctx context.Context, req *rooms.AddToRoomRequest) (*rooms.AddToRoomResponse, error) {
	u, ok := ctx.Value(interceptor.UserContextKey).(*models.User)
	if !ok {
		s.l.Error("Cant get user from context")
		return nil, status.Error(codes.Internal, "Cant authenticate user")
	}

	room, err := s.storage.GetRoom(ctx, req.RoomId)
	if err != nil {
		s.l.Error("Cant get room", slog.String("error", err.Error()))
		return nil, status.Error(codes.NotFound, "Room not found")
	}

	if !room.HasUser(u) {
		return nil, status.Error(codes.NotFound, "User is not in the room")
	}

	newUsers := make([]*models.User, 0, len(room.Users)-1)
	for _, user := range room.Users {
		if user.ID != u.ID {
			newUsers = append(newUsers, user)
		}
	}

	room.Users = newUsers

	roomResp, err := s.storage.UpdateRoom(ctx, room)
	if err != nil {
		s.l.Error("Cant update room", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "Cant update room")
	}

	return &rooms.AddToRoomResponse{Room: roomResp.ToProto()}, nil
}

func (s *Service) GetRoomsByUser(ctx context.Context, req *rooms.GetRoomsByUserRequest) (*rooms.GetRoomsByUserResponse, error) {
	u, ok := ctx.Value(interceptor.UserContextKey).(*models.User)
	if !ok {
		s.l.Error("Cant get user from context")
		return nil, status.Error(codes.Internal, "Cant authenticate user")
	}

	userRooms, err := s.storage.GetRoomsByUser(ctx, u)
	if err != nil {
		s.l.Error("Cant get rooms", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "Cant get rooms")
	}

	roomsProto := make([]*rooms.Room, 0, len(userRooms))
	for _, room := range userRooms {
		roomsProto = append(roomsProto, room.ToProto())
	}

	return &rooms.GetRoomsByUserResponse{Rooms: roomsProto}, nil
}
