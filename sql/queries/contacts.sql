-- name: CreateContact :one
INSERT into contacts (
  id, name, email, message
) VALUES (
$1, $2, $3, $4
)
returning *;
