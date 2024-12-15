// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workouts.sql

package repository

import (
	"context"
)

const countAllWorkouts = `-- name: CountAllWorkouts :one
SELECT COUNT(*) from workouts
`

func (q *Queries) CountAllWorkouts(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countAllWorkouts)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createWorkoutAndReturnId = `-- name: CreateWorkoutAndReturnId :one
INSERT INTO workouts (
  id, name, created_at, updated_at
) VALUES (
  ?1, ?2, ?3, ?4
)
RETURNING id
`

type CreateWorkoutAndReturnIdParams struct {
	ID        string
	Name      string
	CreatedAt string
	UpdatedAt string
}

func (q *Queries) CreateWorkoutAndReturnId(ctx context.Context, arg CreateWorkoutAndReturnIdParams) (string, error) {
	row := q.db.QueryRowContext(ctx, createWorkoutAndReturnId,
		arg.ID,
		arg.Name,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const getAllWorkouts = `-- name: GetAllWorkouts :many
SELECT id, name, created_at, updated_at FROM workouts 
ORDER by id
`

func (q *Queries) GetAllWorkouts(ctx context.Context) ([]Workout, error) {
	rows, err := q.db.QueryContext(ctx, getAllWorkouts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Workout
	for rows.Next() {
		var i Workout
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
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
SELECT id, name, created_at, updated_at FROM workouts 
WHERE id = ?1
`

func (q *Queries) GetWorkoutById(ctx context.Context, id string) (Workout, error) {
	row := q.db.QueryRowContext(ctx, getWorkoutById, id)
	var i Workout
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
