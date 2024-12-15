-- name: CountAllWorkouts :one
SELECT COUNT(*) from workouts;

-- name: GetAllWorkouts :many
SELECT * FROM workouts 
ORDER by id;

-- name: GetWorkoutById :one
SELECT * FROM workouts 
WHERE id = sqlc.arg(id);

-- name: CreateWorkoutAndReturnId :one
INSERT INTO workouts (
  id, name, created_at, updated_at
) VALUES (
  sqlc.arg(id), sqlc.arg(name), sqlc.arg(created_at), sqlc.arg(updated_at)
)
RETURNING id;
