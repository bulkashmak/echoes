-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password_hash)
VALUES (gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2)
RETURNING *;

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

-- name: GetUserFromRefreshToken :one
SELECT u.*
FROM users u
         JOIN refresh_tokens rf ON u.id = rf.user_id
WHERE rf.token = $1;

-- name: UpdateUserEmailAndPasswordByID :one
UPDATE users SET email = $2, password_hash = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateEchoesRed :one
UPDATE users
SET updated_at = NOW(),
    is_echoes_red = true
WHERE id = $1
RETURNING *;

