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

-- name: GetByUserId :one
SELECT * FROM users 
WHERE id = sqlc.arg(id);

-- name: UpdateUser :execrows
UPDATE users
SET username = sqlc.arg(username), password = sqlc.arg(password), updated_on = sqlc.arg(updated_on)
WHERE id = sqlc.arg(id);
