-- name: GetRegistrationUsers :many
SELECT *
FROM users
ORDER BY created_at DESC;

-- name: FindRegistrationUser :one
SELECT id, name, email
FROM users
WHERE email = ?;

-- name: AddUser :exec
INSERT INTO users (id, name, email)
VALUES (?, ?, ?);
