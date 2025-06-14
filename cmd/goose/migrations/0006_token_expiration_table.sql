-- +goose Up
-- +goose StatementBegin
CREATE TABLE expired_tokens (
    id integer primary key autoincrement,
    token text not null,
    token_type text not null,
    created_on text not null,
    remove_on text not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE expired_tokens;
-- +goose StatementEnd
