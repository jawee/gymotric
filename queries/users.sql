-- name: GetByUsername :one
SELECT * FROM users 
WHERE username = sqlc.arg(username);

-- name: CreateUserAndReturnId :one
INSERT INTO users (
  id, username, password, created_on, updated_on
) VALUES (
  sqlc.arg(id), sqlc.arg(username), sqlc.arg(password), sqlc.arg(created_on), sqlc.arg(updated_on)
)
RETURNING id;
