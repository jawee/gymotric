// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: exercise-types.sql

package repository

import (
	"context"
)

const countAllExerciseTypes = `-- name: CountAllExerciseTypes :one
SELECT COUNT(*) from exercise_types
`

func (q *Queries) CountAllExerciseTypes(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countAllExerciseTypes)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createExerciseTypeAndReturnId = `-- name: CreateExerciseTypeAndReturnId :one
INSERT INTO exercise_types (
  id, name
) VALUES (
  ?1, ?2
)
RETURNING id
`

type CreateExerciseTypeAndReturnIdParams struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) CreateExerciseTypeAndReturnId(ctx context.Context, arg CreateExerciseTypeAndReturnIdParams) (string, error) {
	row := q.db.QueryRowContext(ctx, createExerciseTypeAndReturnId, arg.ID, arg.Name)
	var id string
	err := row.Scan(&id)
	return id, err
}

const getAllExerciseTypes = `-- name: GetAllExerciseTypes :many
SELECT id, name FROM exercise_types 
ORDER by id asc
`

func (q *Queries) GetAllExerciseTypes(ctx context.Context) ([]ExerciseType, error) {
	rows, err := q.db.QueryContext(ctx, getAllExerciseTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExerciseType
	for rows.Next() {
		var i ExerciseType
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const getExerciseTypeById = `-- name: GetExerciseTypeById :one
SELECT id, name FROM exercise_types 
WHERE id = ?1
`

func (q *Queries) GetExerciseTypeById(ctx context.Context, id string) (ExerciseType, error) {
	row := q.db.QueryRowContext(ctx, getExerciseTypeById, id)
	var i ExerciseType
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}
