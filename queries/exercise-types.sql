-- name: CountAllExerciseTypes :one
SELECT COUNT(*) from exercise_types;

-- name: GetAllExerciseTypes :many
SELECT * FROM exercise_types 
ORDER by id asc;

-- name: GetExerciseTypeById :one
SELECT * FROM exercise_types 
WHERE id = sqlc.arg(id);
