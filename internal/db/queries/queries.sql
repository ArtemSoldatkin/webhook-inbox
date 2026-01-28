-- name: RegisterEndpoint :one
INSERT INTO endpoints (
    user_id,
    url,
    name,
    description,
    headers
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;
