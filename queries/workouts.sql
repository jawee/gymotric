-- name: CountAllWorkouts :one
SELECT COUNT(*) from Workouts;

-- name: CountVisitors :one
SELECT COUNT(distinct(ip)) FROM visit WHERE visited_at > $1;

-- name: InsertVisit :exec
INSERT INTO visit (ip, visited_at) VALUES ((sqlc.arg(ip)::varchar)::inet, now());
