-- name: CreateUserDir :one
INSERT INTO directory (
  name,
  user_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUserDirs :many
SELECT * FROM directory
WHERE user_id = $1
ORDER BY dir_id;

-- name: CheckUserDir :one
SELECT * FROM directory
WHERE user_id = $1 AND name = $2
LIMIT 1;

-- name: DeleteUserDir :exec
DELETE FROM directory
WHERE user_id = $1 AND name = $2;
