-- name: GetAllExerciseTypes :many
SELECT * FROM exercise_types 
ORDER by id asc;

-- name: GetExerciseTypeById :one
SELECT * FROM exercise_types 
WHERE id = sqlc.arg(id);

-- name: CreateExerciseTypeAndReturnId :one
INSERT INTO exercise_types (
  id, name, created_on, updated_on
) VALUES (
  sqlc.arg(id), sqlc.arg(name), sqlc.arg(created_on), sqlc.arg(updated_on)
)
RETURNING id;

-- name: DeleteExerciseTypeById :execrows
DELETE FROM exercise_types
WHERE id = sqlc.arg(id);
