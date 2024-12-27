package model

import "time"

type User struct {
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
