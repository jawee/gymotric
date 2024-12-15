-- name: CountAllSets :one
SELECT COUNT(*) from sets;

-- name: GetAllSets :many
SELECT * FROM sets 
ORDER by id;

-- name: GetSetById :one
SELECT * FROM sets 
WHERE id = sqlc.arg(id);

-- name: CreateSetAndReturnId :one
INSERT INTO sets (
  id, repetitions, weight, exercise_id
) VALUES (
  sqlc.arg(id), sqlc.arg(repetitions), sqlc.arg(weight), sqlc.arg(exercise_id)
)
RETURNING id;

-- name: GetSetsByExerciseId :many
SELECT * FROM sets 
WHERE exercise_id = sqlc.arg(exercise_id);
