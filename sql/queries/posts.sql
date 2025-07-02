-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2) RETURNING *;

-- name: ListPosts :many
SELECT *
FROM posts
ORDER BY posts.created_at;

-- name: RetrievePostByID :one
SELECT *
FROM posts
WHERE posts.id = $1;

-- name: DeletePostByID :exec
DELETE FROM posts
WHERE posts.id = $1;

