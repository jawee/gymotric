-- name: GetAllWorkouts :many
SELECT * FROM workouts 
ORDER by id;

-- name: GetWorkoutById :one
SELECT * FROM workouts 
WHERE id = sqlc.arg(id);

-- name: CreateWorkoutAndReturnId :one
INSERT INTO workouts (
  id, name, created_on, updated_on
) VALUES (
  sqlc.arg(id), sqlc.arg(name), sqlc.arg(created_on), sqlc.arg(updated_on)
)
RETURNING id;

-- name: CompleteWorkoutById :exec
UPDATE workouts 
set completed_on = sqlc.arg(completed_on) 
where id = sqlc.arg(id)
