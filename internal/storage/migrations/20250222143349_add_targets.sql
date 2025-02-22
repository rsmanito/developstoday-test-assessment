-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS targets (
  id SERIAL PRIMARY KEY,
  mission INT NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
  name VARCHAR(30) NOT NULL,
  country VARCHAR(30) NOT NULL,
  notes VARCHAR(256) NOT NULL,
  completed BOOLEAN NOT NULL DEFAULT false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS targets;
-- +goose StatementEnd
