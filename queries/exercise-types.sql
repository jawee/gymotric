-- name: GetAllExerciseTypes :many
SELECT * FROM exercise_types 
ORDER by id asc;

-- name: GetExerciseTypeById :one
SELECT * FROM exercise_types 
WHERE id = sqlc.arg(id);

-- name: CreateExerciseTypeAndReturnId :one
INSERT INTO exercise_types (
  id, name
) VALUES (
  sqlc.arg(id), sqlc.arg(name)
)
RETURNING id;
