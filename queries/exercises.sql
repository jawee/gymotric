-- name: GetAllExercises :many
SELECT * FROM exercises 
ORDER by id;

-- name: GetExerciseById :one
SELECT * FROM exercises 
WHERE id = sqlc.arg(id);

-- name: CreateExerciseAndReturnId :one
INSERT INTO exercises (
  id, name, workout_id, exercise_type_id
) VALUES (
  sqlc.arg(id), sqlc.arg(name), sqlc.arg(workout_id), sqlc.arg(exercise_type_id)
)
RETURNING id;

-- name: GetExercisesByWorkoutId :many
SELECT * FROM exercises
WHERE workout_id = sqlc.arg(workout_id)
