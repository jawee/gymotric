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
SELECT COALESCE(s.weight, 0) as weight, COALESCE(Max(s.repetitions), 0) as repetitions FROM exercises e
LEFT JOIN sets s ON s.exercise_id = e.id
WHERE e.exercise_type_id = sqlc.arg(id) AND (s.user_id = sqlc.arg(user_id) OR s.id IS NULL)
AND (s.weight = (SELECT Max(s.weight) as weight FROM exercises e
JOIN sets s ON s.exercise_id = e.id
WHERE e.exercise_type_id = sqlc.arg(id) AND s.user_id = sqlc.arg(user_id)) OR s.weight IS NULL)
GROUP BY e.exercise_type_id;

-- name: GetLastWeightRepsByExerciseTypeId :one
SELECT s.repetitions, s.weight FROM exercises e
JOIN sets s ON s.exercise_id = e.id
JOIN workouts w ON e.workout_id = w.id
WHERE exercise_type_id = sqlc.arg(id) 
AND s.user_id = sqlc.arg(user_id)
AND w.completed_on IS NOT NULL
ORDER BY s.id desc LIMIT 1;


-- name: UpdateExerciseType :execrows
UPDATE exercise_types
SET name = sqlc.arg(name), updated_on = sqlc.arg(updated_on)
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);
