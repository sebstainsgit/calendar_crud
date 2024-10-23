// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	EventID   uuid.UUID
	EventName string
	UsersID   uuid.UUID
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Login struct {
	RefrToken string
	UsersID   uuid.UUID
	Expires   time.Time
}

type User struct {
	UserID    uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}