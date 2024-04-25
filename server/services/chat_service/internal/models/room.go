package models

type Room struct {
	ID        string
	Name      string
	UserIDS   []string
	CreatedBy string
}
