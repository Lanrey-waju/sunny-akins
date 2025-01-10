-- +goose Up
CREATE TABLE users
(
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY, 
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR NOT NULL,
  email VARCHAR(255) NOT NULL,
  hashed_password CHAR(60) NOT NULL
);

-- +goose Down
DROP TABLE users;
