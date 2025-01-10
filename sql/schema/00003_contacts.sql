-- +goose Up
-- +goose StatementBegin
CREATE TABLE contacts (
  id UUID default gen_random_uuid() primary key,
  name varchar(50) not null,
  email varchar(255) not null,
  message text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE contacts;
-- +goose StatementEnd
