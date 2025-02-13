-- name: AddExample :one
INSERT INTO example (id) VALUES ($1) RETURNING *;

-- name: GetExample :one
SELECT * FROM example WHERE id = $1;
