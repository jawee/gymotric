-- name: GetAllSets :many
SELECT * FROM sets 
ORDER by id;

-- name: GetSetById :one
SELECT * FROM sets 
WHERE id = sqlc.arg(id);

-- name: CreateSetAndReturnId :one
INSERT INTO sets (
  id, repetitions, weight, exercise_id, created_on, updated_on
) VALUES (
  sqlc.arg(id), sqlc.arg(repetitions), sqlc.arg(weight), sqlc.arg(exercise_id), sqlc.arg(created_on), sqlc.arg(updated_on)
)
RETURNING id;

-- name: GetSetsByExerciseId :many
SELECT * FROM sets 
WHERE exercise_id = sqlc.arg(exercise_id);

-- name: DeleteSetById :execrows
DELETE FROM sets 
where id = sqlc.arg(id);
