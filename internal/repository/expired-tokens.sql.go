// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: expired-tokens.sql

package repository

import (
	"context"
)

const checkIfTokenExists = `-- name: CheckIfTokenExists :one
SELECT id FROM expired_tokens
WHERE token = ?1 
AND token_type = ?2
`

type CheckIfTokenExistsParams struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

func (q *Queries) CheckIfTokenExists(ctx context.Context, arg CheckIfTokenExistsParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkIfTokenExists, arg.Token, arg.TokenType)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createExpiredToken = `-- name: CreateExpiredToken :execrows
INSERT INTO expired_tokens (
  token, token_type, created_on, remove_on
) VALUES (
  ?1, ?2, ?3, ?4
)
`

type CreateExpiredTokenParams struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	CreatedOn string `json:"created_on"`
	RemoveOn  string `json:"remove_on"`
}

func (q *Queries) CreateExpiredToken(ctx context.Context, arg CreateExpiredTokenParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createExpiredToken,
		arg.Token,
		arg.TokenType,
		arg.CreatedOn,
		arg.RemoveOn,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const deleteExpiredTokens = `-- name: DeleteExpiredTokens :execrows
DELETE FROM expired_tokens
WHERE remove_on < ?1
`

func (q *Queries) DeleteExpiredTokens(ctx context.Context, currTime string) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteExpiredTokens, currTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
