-- name: GetAllExercises :many
SELECT * FROM exercises 
ORDER by id;

-- name: GetExerciseById :one
SELECT * FROM exercises 
WHERE id = sqlc.arg(id);

-- name: CreateExerciseAndReturnId :one
INSERT INTO exercises (
  id, name, workout_id, exercise_type_id, created_on, updated_on
) VALUES (
  sqlc.arg(id), sqlc.arg(name), sqlc.arg(workout_id), sqlc.arg(exercise_type_id), sqlc.arg(created_on), sqlc.arg(updated_on)
)
RETURNING id;

-- name: GetExercisesByWorkoutId :many
SELECT * FROM exercises
WHERE workout_id = sqlc.arg(workout_id);

-- name: DeleteExerciseById :execrows
DELETE FROM exercises where id = sqlc.arg(id);
