-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password_hash)
VALUES (gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2) RETURNING *;

-- name: DeleteAllUsers :exec
DELETE
FROM users;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE users.id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE users.email = $1;
