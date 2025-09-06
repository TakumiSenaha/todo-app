-- name: CreateTodo :one
INSERT INTO todos (
    user_id,
    title,
    due_date,
    priority,
    is_completed
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1 LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todos
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateTodo :one
UPDATE todos
SET title = $2,
    due_date = $3,
    priority = $4,
    is_completed = $5
WHERE id = $1 AND user_id = $6
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1 AND user_id = $2;
