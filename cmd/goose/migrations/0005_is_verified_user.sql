-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN is_verified boolean not null DEFAULT false;

UPDATE users
SET is_verified = true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN is_verified;
-- +goose StatementEnd
