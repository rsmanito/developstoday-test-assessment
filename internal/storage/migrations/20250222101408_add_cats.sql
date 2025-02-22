-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cats (
  id SERIAL PRIMARY KEY,
  name VARCHAR(30) NOT NULL,
  years_of_experience INT NOT NULL DEFAULT 0,
  breed VARCHAR(30) NOT NULL,
  salary INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cats;
-- +goose StatementEnd
