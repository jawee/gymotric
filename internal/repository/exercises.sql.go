// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: exercises.sql

package repository

import (
	"context"
)

const createExerciseAndReturnId = `-- name: CreateExerciseAndReturnId :one
INSERT INTO exercises (
  id, name, workout_id, exercise_type_id
) VALUES (
  ?1, ?2, ?3, ?4
)
RETURNING id
`

type CreateExerciseAndReturnIdParams struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	WorkoutID      string `json:"workout_id"`
	ExerciseTypeID string `json:"exercise_type_id"`
}

func (q *Queries) CreateExerciseAndReturnId(ctx context.Context, arg CreateExerciseAndReturnIdParams) (string, error) {
	row := q.db.QueryRowContext(ctx, createExerciseAndReturnId,
		arg.ID,
		arg.Name,
		arg.WorkoutID,
		arg.ExerciseTypeID,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const getAllExercises = `-- name: GetAllExercises :many
SELECT id, name, workout_id, exercise_type_id FROM exercises 
ORDER by id
`

func (q *Queries) GetAllExercises(ctx context.Context) ([]Exercise, error) {
	rows, err := q.db.QueryContext(ctx, getAllExercises)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Exercise{}
	for rows.Next() {
		var i Exercise
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.WorkoutID,
			&i.ExerciseTypeID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getExerciseById = `-- name: GetExerciseById :one
SELECT id, name, workout_id, exercise_type_id FROM exercises 
WHERE id = ?1
`

func (q *Queries) GetExerciseById(ctx context.Context, id string) (Exercise, error) {
	row := q.db.QueryRowContext(ctx, getExerciseById, id)
	var i Exercise
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.WorkoutID,
		&i.ExerciseTypeID,
	)
	return i, err
}

const getExercisesByWorkoutId = `-- name: GetExercisesByWorkoutId :many
SELECT id, name, workout_id, exercise_type_id FROM exercises
WHERE workout_id = ?1
`

func (q *Queries) GetExercisesByWorkoutId(ctx context.Context, workoutID string) ([]Exercise, error) {
	rows, err := q.db.QueryContext(ctx, getExercisesByWorkoutId, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Exercise{}
	for rows.Next() {
		var i Exercise
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.WorkoutID,
			&i.ExerciseTypeID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
