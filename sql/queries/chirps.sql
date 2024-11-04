-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2) RETURNING *;

-- name: ListChirps :many
SELECT id, created_at, updated_at, body, user_id
FROM chirps
ORDER BY created_at;

-- name: ListChirpsByAuthor :many
SELECT id, created_at, updated_at, body, user_id
FROM chirps
WHERE user_id = $1
ORDER BY created_at;

-- name: GetChirpById :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteChirpById :exec
DELETE FROM chirps
WHERE id = $1;