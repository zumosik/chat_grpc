package models

import (
	"github.com/zumosik/grpc_chat_protos/gen/go/auth"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                string `db:"id"`
	Username          string `db:"username"`
	Password          string `db:"password"`
	Email             string `db:"email"`
	EncryptedPassword []byte `db:"encrypted_password"`
}

// ToAuthUser converts the *User to an *auth.User
func (u *User) ToAuthUser() *auth.User {
	return &auth.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

// HashPassword hashes the password of the user and stores it in the EncryptedPassword field
func (u *User) HashPassword() error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.EncryptedPassword = encryptedPassword

	return nil
}

// ComparePassword compares the password of the user with the provided password
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.EncryptedPassword, []byte(password)) == nil
}
