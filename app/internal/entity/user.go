package entity

import "time"

type User struct {
	UserId       string
	UserType     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
