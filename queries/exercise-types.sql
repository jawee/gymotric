-- name: GetAllExerciseTypes :many
SELECT * FROM exercise_types 
WHERE user_id = sqlc.arg(user_id)
ORDER by id asc;

-- name: GetExerciseTypeById :one
SELECT * FROM exercise_types 
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);

-- name: CreateExerciseTypeAndReturnId :one
INSERT INTO exercise_types (
  id, name, created_on, updated_on, user_id
) VALUES (
  sqlc.arg(id), sqlc.arg(name), sqlc.arg(created_on), sqlc.arg(updated_on), sqlc.arg(user_id)
)
RETURNING id;

-- name: DeleteExerciseTypeById :execrows
DELETE FROM exercise_types
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);


-- name: GetMaxWeightRepsByExerciseTypeId :one
SELECT s.repetitions, Max(s.weight) FROM exercises e
JOIN sets s ON s.exercise_id = e.id
WHERE exercise_type_id = sqlc.arg(id) AND s.user_id = sqlc.arg(user_id);

-- name: GetLastWeightRepsByExerciseTypeId :one
SELECT s.repetitions, s.weight FROM exercises e
JOIN sets s ON s.exercise_id = e.id
JOIN workouts w ON e.workout_id = w.id
WHERE exercise_type_id = sqlc.arg(id) 
AND s.user_id = sqlc.arg(user_id)
AND w.completed_on IS NOT NULL
ORDER BY s.id desc LIMIT 1;

