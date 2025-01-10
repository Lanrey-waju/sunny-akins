// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID      uuid.UUID
	Name    string
	Email   string
	Message string
}

type Session struct {
	Token  string
	Data   []byte
	Expiry time.Time
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Email          string
	HashedPassword string
}
