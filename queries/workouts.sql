-- name: GetAllWorkouts :many
SELECT * FROM workouts 
WHERE user_id = sqlc.arg(user_id)
ORDER BY id DESC
LIMIT sqlc.arg(limit) OFFSET sqlc.arg(offset);

-- name: GetAllWorkoutsCount :one
SELECT count(*) FROM workouts 
WHERE user_id = sqlc.arg(user_id);

-- name: GetWorkoutById :one
SELECT * FROM workouts 
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);

-- name: CreateWorkoutAndReturnId :one
INSERT INTO workouts (
  id, name, created_on, updated_on, user_id
) VALUES (
  sqlc.arg(id), sqlc.arg(name), sqlc.arg(created_on), sqlc.arg(updated_on), sqlc.arg(user_id)
)
RETURNING id;

-- name: CompleteWorkoutById :execrows
UPDATE workouts 
SET completed_on = sqlc.arg(completed_on), updated_on = sqlc.arg(updated_on)
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);

-- name: DeleteWorkoutById :execrows
DELETE FROM workouts
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);

-- name: UpdateWorkoutById :execrows
UPDATE workouts
SET note = sqlc.arg(note), updated_on = sqlc.arg(updated_on)
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);

-- name: ReopenWorkoutById :execrows
UPDATE workouts
SET completed_on = NULL, updated_on = sqlc.arg(updated_on)
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);
