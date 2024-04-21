package models

import "github.com/zumosik/grpc_chat_protos/go/rooms"

type User struct {
	ID       string
	Username string
	Email    string
	Verified bool
}

type Room struct {
	ID        string
	Name      string
	Users     []*User
	CreatedBy *User
}

func (r *Room) ToProto() *rooms.Room {
	users := make([]*rooms.PublicUser, 0, len(r.Users))
	for _, user := range r.Users {
		users = append(users, &rooms.PublicUser{
			Username: user.Username,
			Email:    user.Email,
			Verified: user.Verified,
		})
	}

	return &rooms.Room{
		Id:    r.ID,
		Name:  r.Name,
		Users: users,
		CreatedBy: &rooms.PublicUser{
			Username: r.CreatedBy.Username,
			Email:    r.CreatedBy.Email,
			Verified: r.CreatedBy.Verified,
		},
	}
}

func (r *Room) HasUser(u *User) bool {
	for _, user := range r.Users {
		if user.ID == u.ID {
			return true
		}
	}
	return false
}
