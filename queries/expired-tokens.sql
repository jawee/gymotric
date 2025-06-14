-- name: CheckIfTokenExists :one
SELECT id FROM expired_tokens
WHERE token = sqlc.arg(token) 
AND token_type = sqlc.arg(token_type);

-- name: CreateExpiredToken :execrows
INSERT INTO expired_tokens (
  token, token_type, created_on, remove_on
) VALUES (
  sqlc.arg(token), sqlc.arg(token_type), sqlc.arg(created_on), sqlc.arg(remove_on)
);

-- name: DeleteExpiredTokens :execrows
DELETE FROM expired_tokens
WHERE remove_on < sqlc.arg(curr_time);
