// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: sets.sql

package repository

import (
	"context"
)

const countAllSets = `-- name: CountAllSets :one
SELECT COUNT(*) from sets
`

func (q *Queries) CountAllSets(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countAllSets)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createSetAndReturnId = `-- name: CreateSetAndReturnId :one
INSERT INTO sets (
  id, repetitions, weight, exercise_id
) VALUES (
  ?1, ?2, ?3, ?4
)
RETURNING id
`

type CreateSetAndReturnIdParams struct {
	ID          string
	Repetitions int64
	Weight      float64
	ExerciseID  string
}

func (q *Queries) CreateSetAndReturnId(ctx context.Context, arg CreateSetAndReturnIdParams) (string, error) {
	row := q.db.QueryRowContext(ctx, createSetAndReturnId,
		arg.ID,
		arg.Repetitions,
		arg.Weight,
		arg.ExerciseID,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const getAllSets = `-- name: GetAllSets :many
SELECT id, repetitions, weight, exercise_id FROM sets 
ORDER by id
`

func (q *Queries) GetAllSets(ctx context.Context) ([]Set, error) {
	rows, err := q.db.QueryContext(ctx, getAllSets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Set
	for rows.Next() {
		var i Set
		if err := rows.Scan(
			&i.ID,
			&i.Repetitions,
			&i.Weight,
			&i.ExerciseID,
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

const getSetById = `-- name: GetSetById :one
SELECT id, repetitions, weight, exercise_id FROM sets 
WHERE id = ?1
`

func (q *Queries) GetSetById(ctx context.Context, id string) (Set, error) {
	row := q.db.QueryRowContext(ctx, getSetById, id)
	var i Set
	err := row.Scan(
		&i.ID,
		&i.Repetitions,
		&i.Weight,
		&i.ExerciseID,
	)
	return i, err
}
