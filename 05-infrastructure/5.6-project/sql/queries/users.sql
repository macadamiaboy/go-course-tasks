-- name: CreateUser :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;