-- name: CreateUser :one
INSERT INTO users (
  username, 
  password, 
  name,
  surname,
  enabled,
  role,
  email
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;


-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE username = $1;

-- name: UpdateUserEmail :one
UPDATE users 
SET email = $1 
WHERE username = $2
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: UpdateUser :one
UPDATE users
SET
  is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
WHERE
  username = sqlc.arg(username)
RETURNING *;