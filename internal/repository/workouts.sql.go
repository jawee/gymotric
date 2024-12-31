// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workouts.sql

package repository

import (
	"context"
)

const completeWorkoutById = `-- name: CompleteWorkoutById :execrows
UPDATE workouts 
set completed_on = ?1 
where id = ?2
`

type CompleteWorkoutByIdParams struct {
	CompletedOn interface{} `json:"completed_on"`
	ID          string      `json:"id"`
}

func (q *Queries) CompleteWorkoutById(ctx context.Context, arg CompleteWorkoutByIdParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, completeWorkoutById, arg.CompletedOn, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const createWorkoutAndReturnId = `-- name: CreateWorkoutAndReturnId :one
INSERT INTO workouts (
  id, name, created_on, updated_on
) VALUES (
  ?1, ?2, ?3, ?4
)
RETURNING id
`

type CreateWorkoutAndReturnIdParams struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

func (q *Queries) CreateWorkoutAndReturnId(ctx context.Context, arg CreateWorkoutAndReturnIdParams) (string, error) {
	row := q.db.QueryRowContext(ctx, createWorkoutAndReturnId,
		arg.ID,
		arg.Name,
		arg.CreatedOn,
		arg.UpdatedOn,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const getAllWorkouts = `-- name: GetAllWorkouts :many
SELECT id, name, completed_on, created_on, updated_on FROM workouts 
ORDER by id
`

func (q *Queries) GetAllWorkouts(ctx context.Context) ([]Workout, error) {
	rows, err := q.db.QueryContext(ctx, getAllWorkouts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Workout{}
	for rows.Next() {
		var i Workout
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CompletedOn,
			&i.CreatedOn,
			&i.UpdatedOn,
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

const getWorkoutById = `-- name: GetWorkoutById :one
SELECT id, name, completed_on, created_on, updated_on FROM workouts 
WHERE id = ?1
`

func (q *Queries) GetWorkoutById(ctx context.Context, id string) (Workout, error) {
	row := q.db.QueryRowContext(ctx, getWorkoutById, id)
	var i Workout
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CompletedOn,
		&i.CreatedOn,
		&i.UpdatedOn,
	)
	return i, err
}
