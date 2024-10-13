-- name: AddExample :one
INSERT INTO example (id, field)  VALUES (?, ?) RETURNING *;

-- name: GetExample :one
SELECT * FROM example WHERE id = ?;

-- name: UpdateExampleField :one
UPDATE example SET field = ? WHERE id = ? RETURNING *;
