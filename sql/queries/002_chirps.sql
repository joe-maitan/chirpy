-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteChirps :one
DELETE FROM chirps RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at;

-- name: GetChirp :one
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2;

-- name: GetAllChripsByUser :many
SELECT * FROM chirps 
WHERE user_id = $1
ORDER BY created_at;