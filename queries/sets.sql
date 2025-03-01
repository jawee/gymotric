-- name: GetAllSets :many
SELECT * FROM sets 
WHERE user_id = sqlc.arg(user_id)
ORDER by id;

-- name: GetSetById :one
SELECT * FROM sets 
WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- name: CreateSetAndReturnId :one
INSERT INTO sets (
  id, repetitions, weight, exercise_id, created_on, updated_on, user_id
) VALUES (
  sqlc.arg(id), sqlc.arg(repetitions), sqlc.arg(weight), sqlc.arg(exercise_id), sqlc.arg(created_on), sqlc.arg(updated_on), sqlc.arg(user_id)
)
RETURNING id;

-- name: GetSetsByExerciseId :many
SELECT * FROM sets 
WHERE exercise_id = sqlc.arg(exercise_id)
AND user_id = sqlc.arg(user_id);

-- name: DeleteSetById :execrows
DELETE FROM sets 
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);
