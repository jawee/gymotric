-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN email text null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN email;
-- +goose StatementEnd
