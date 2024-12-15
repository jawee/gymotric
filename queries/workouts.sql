-- name: CountAllWorkouts :one
SELECT COUNT(*) from workouts;

-- name: GetAllWorkouts :many
SELECT * FROM workouts 
ORDER by id;

-- name: GetWorkoutById :one
SELECT * FROM workouts 
WHERE id = sqlc.arg(id);
