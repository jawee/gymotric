-- name: GetStatisticsSinceDate :one
SELECT count(*) FROM
workouts
WHERE user_id = sqlc.arg(user_id) AND
completed_on >= sqlc.arg(start_date);

-- name: GetStatisticsBetweenDates :one
SELECT count(*) FROM
workouts
WHERE user_id = sqlc.arg(user_id) AND
completed_on >= sqlc.arg(start_date) AND
completed_on <= sqlc.arg(end_date);
