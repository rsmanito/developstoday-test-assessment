-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS missions (
  id SERIAL PRIMARY KEY,
  assignee INT REFERENCES cats(id) ON DELETE SET NULL DEFAULT NULL,
  completed BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS missions;
-- +goose StatementEnd
