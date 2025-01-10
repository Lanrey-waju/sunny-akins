// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: contacts.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createContact = `-- name: CreateContact :one
INSERT into contacts (
  id, name, email, message
) VALUES (
$1, $2, $3, $4
)
returning id, name, email, message
`

type CreateContactParams struct {
	ID      uuid.UUID
	Name    string
	Email   string
	Message string
}

func (q *Queries) CreateContact(ctx context.Context, arg CreateContactParams) (Contact, error) {
	row := q.db.QueryRowContext(ctx, createContact,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Message,
	)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Message,
	)
	return i, err
}
