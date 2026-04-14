-- name: CreateUser :one
INSERT INTO users (email) VALUES ($1) RETURNING id, email;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id;