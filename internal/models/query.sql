-- name: GetTask :one
SELECT *
FROM scheduler
WHERE id = ?
LIMIT 1;
-- name: ListTasks :many
SELECT *
FROM scheduler
ORDER BY date;
-- name: CreateTask :one
INSERT INTO scheduler (date, title, comment, repeat)
VALUES (?, ?, ?, ?)
RETURNING *;
-- name: UpdateTask :exec
UPDATE scheduler
set date = ?,
    title = ?,
    comment = ?,
    repeat = ?
WHERE id = ?;
-- name: DeleteTask :exec
DELETE FROM scheduler
WHERE id = ?;