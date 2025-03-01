-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id text primary key,
    username text unique not null,
    password text not null,
    created_on text not null,
    updated_on text not null
);

CREATE TABLE workouts (
    id text primary key,
    name text not null,
    completed_on text null,

    created_on text not null,
    updated_on text not null,

    user_id text not null,

    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE exercise_types (
    id text primary key,
    name text not null,

    created_on text not null,
    updated_on text not null,

    user_id text not null,

    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE exercises (
    id text primary key,
    name text not null,
    created_on text not null,
    updated_on text not null,

    user_id text not null,
    workout_id text not null,
    exercise_type_id text not null,

    FOREIGN KEY(user_id) REFERENCES users(id)
    FOREIGN KEY(exercise_type_id) REFERENCES exercise_types(id),
    FOREIGN KEY(workout_id) REFERENCES workouts(id)

);

CREATE TABLE sets (
    id text primary key,
    repetitions integer not null,
    weight real not null,
    created_on text not null,
    updated_on text not null,

    user_id text not null,
    exercise_id text not null,

    FOREIGN KEY(user_id) REFERENCES users(id)
    FOREIGN KEY(exercise_id) REFERENCES exercises(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sets;
DROP TABLE exercises;
DROP TABLE exercise_types;
DROP TABLE workouts;
DROP TABLE users;
-- +goose StatementEnd
