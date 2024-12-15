-- name: CountAllExercises :one
SELECT COUNT(*) from exercises;

-- name: GetAllExercises :many
SELECT * FROM exercises 
ORDER by id;

-- name: GetExerciseById :one
SELECT * FROM exercises 
WHERE id = sqlc.arg(id);
