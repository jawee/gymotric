-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX unique_email ON users(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS unique_email;
-- +goose StatementEnd
