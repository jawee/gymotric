-- +goose Up
-- +goose StatementBegin
CREATE TABLE workouts (
    id text primary key,
    name text not null,

    created_at text not null,
    updated_at text not null
);

CREATE TABLE exercises (
    id text primary key,
    name text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workouts;
DROP TABLE exercises;
-- +goose StatementEnd
