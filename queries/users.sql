-- name: GetByUsername :one
SELECT * FROM users 
WHERE username = sqlc.arg(username);

-- name: CreateUserAndReturnId :one
INSERT INTO users (
  id, username, password, created_on, updated_on, email
) VALUES (
  sqlc.arg(id), sqlc.arg(username), sqlc.arg(password), sqlc.arg(created_on), sqlc.arg(updated_on), sqlc.arg(email)
)
RETURNING id;

-- name: GetByUserId :one
SELECT * FROM users 
WHERE id = sqlc.arg(id);

-- name: UpdateUser :execrows
UPDATE users
SET password = sqlc.arg(password), updated_on = sqlc.arg(updated_on), email = sqlc.arg(email), is_verified = sqlc.arg(is_verified)
WHERE id = sqlc.arg(id);

-- name: EmailExists :one
SELECT count(*) from Users
WHERE email = sqlc.arg(email);

-- name: GetByEmail :one
SELECT * FROM users 
WHERE email = sqlc.arg(email);
