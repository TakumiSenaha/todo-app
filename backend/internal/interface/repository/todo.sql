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

-- 完了切り替え専用クエリ
-- name: ToggleTodoComplete :one
UPDATE todos
SET is_completed = NOT is_completed,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2
RETURNING *;

-- ソート機能付きリスト取得
-- name: ListTodosWithSort :many
SELECT * FROM todos
WHERE user_id = $1
ORDER BY
    CASE WHEN $2 = 'due_date_asc' THEN due_date END ASC,
    CASE WHEN $2 = 'due_date_desc' THEN due_date END DESC,
    CASE WHEN $2 = 'priority_desc' THEN priority END DESC,
    CASE WHEN $2 = 'created_desc' THEN created_at END DESC,
    is_completed ASC,
    created_at DESC;
