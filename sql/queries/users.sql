-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    FALSE
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: EmailLookup :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET
    email = $2,
    hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: UpgradeToRed :one
UPDATE users
SET
    is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;

-- name: DowngradeFromRed :one
UPDATE users
SET
    is_chirpy_red = FALSE
WHERE id = $1
RETURNING *;