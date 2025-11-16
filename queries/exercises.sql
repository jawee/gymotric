-- name: GetAllExercises :many
SELECT * FROM exercises 
WHERE user_id = sqlc.arg(user_id)
ORDER by id;

-- name: GetExerciseById :one
SELECT * FROM exercises 
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);

-- name: CreateExerciseAndReturnId :one
INSERT INTO exercises (
  id, name, workout_id, exercise_type_id, exercise_item_id, created_on, updated_on, user_id
) VALUES (
  sqlc.arg(id), sqlc.arg(name), sqlc.arg(workout_id), sqlc.arg(exercise_type_id), sqlc.arg(exercise_item_id), sqlc.arg(created_on), sqlc.arg(updated_on), sqlc.arg(user_id)
)
RETURNING id;

-- name: GetExercisesByWorkoutId :many
SELECT * FROM exercises
WHERE workout_id = sqlc.arg(workout_id)
AND user_id = sqlc.arg(user_id);

-- name: GetExercisesByExerciseItemId :many
SELECT * FROM exercises
WHERE exercise_item_id = sqlc.arg(exercise_item_id)
AND user_id = sqlc.arg(user_id);

-- name: DeleteExerciseById :execrows
DELETE FROM exercises 
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id);
