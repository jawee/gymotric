-- +goose Up
-- +goose StatementBegin

ALTER TABLE exercises RENAME TO old_exercises;

CREATE TABLE exercises (
    id text primary key,
    name text not null,
    created_on text not null,
    updated_on text not null,

    user_id text not null,
    workout_id text not null,
    exercise_type_id text not null,
    exercise_item_id text not null,

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(exercise_type_id) REFERENCES exercise_types(id),
    FOREIGN KEY(workout_id) REFERENCES workouts(id),
    FOREIGN KEY(exercise_item_id) REFERENCES exercise_items(id)
);

INSERT INTO exercises (
    id,
    name,
    created_on,
    updated_on,
    user_id,
    workout_id,
    exercise_type_id,
    exercise_item_id
)
SELECT
    id,
    name,
    created_on,
    updated_on,
    user_id,
    workout_id,
    exercise_type_id,
    exercise_item_id
FROM old_exercises;

DROP TABLE old_exercises;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE exercises RENAME TO old_exercises;

CREATE TABLE exercises (
    id text primary key,
    name text not null,
    created_on text not null,
    updated_on text not null,

    user_id text not null,
    workout_id text not null,
    exercise_type_id text not null,
    exercise_item_id text,

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(exercise_type_id) REFERENCES exercise_types(id),
    FOREIGN KEY(workout_id) REFERENCES workouts(id),
    FOREIGN KEY(exercise_item_id) REFERENCES exercise_items(id)
);

INSERT INTO exercises (
    id,
    name,
    created_on,
    updated_on,
    user_id,
    workout_id,
    exercise_type_id,
    exercise_item_id
)
SELECT
    id,
    name,
    created_on,
    updated_on,
    user_id,
    workout_id,
    exercise_type_id,
    exercise_item_id
FROM old_exercises;

DROP TABLE old_exercises;
-- +goose StatementEnd
