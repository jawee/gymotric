-- name: CreateExerciseItemAndReturnId :one
INSERT INTO exercise_items (
  id, type, user_id, workout_id, created_on, updated_on
) VALUES (
  sqlc.arg(id), sqlc.arg(type), sqlc.arg(user_id), sqlc.arg(workout_id), sqlc.arg(created_on), sqlc.arg(updated_on)
)
RETURNING id;

-- name: GetExerciseItemById :one
SELECT * FROM exercise_items
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);

-- name: GetExerciseItemsByWorkoutId :many
SELECT * FROM exercise_items
WHERE workout_id = sqlc.arg(workout_id)
AND user_id = sqlc.arg(user_id)
ORDER BY created_on;

-- name: UpdateExerciseItemType :execrows
UPDATE exercise_items
SET type = sqlc.arg(type), updated_on = sqlc.arg(updated_on)
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);

-- name: DeleteExerciseItemById :execrows
DELETE FROM exercise_items
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);
