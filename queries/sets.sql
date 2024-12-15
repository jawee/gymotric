-- name: CountAllSets :one
SELECT COUNT(*) from sets;

-- name: GetAllSets :many
SELECT * FROM sets 
ORDER by id;

-- name: GetSetById :one
SELECT * FROM sets 
WHERE id = sqlc.arg(id);
