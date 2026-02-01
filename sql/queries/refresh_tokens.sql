-- name: CreateToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    NOW() + INTERVAL '60 days',
    $2
)
RETURNING *;

-- name: RevokeToken :one
UPDATE refresh_tokens
SET 
    updated_at = NOW(),
    revoked_at = NOW()
WHERE token = $1
RETURNING *;

-- name: GetToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL;
