-- name: RegisterEndpoint :one
INSERT INTO endpoints (
    user_id,
    public_key,
    name,
    description
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;
