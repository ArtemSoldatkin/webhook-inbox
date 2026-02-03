-- name: CreateUser :one
INSERT INTO users (
    email
)
VALUES (
    $1
)
RETURNING *;

-- name: ListEndpoints :many
SELECT
    id,
    user_id,
    url,
    name,
    description,
    headers,
    is_active,
    created_at
FROM
    endpoints
WHERE
    user_id = $1
ORDER BY
    created_at DESC;


-- name: RegisterEndpoint :one
INSERT INTO endpoints (
    user_id,
    url,
    name,
    description,
    headers
)
VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;


-- name: ToggleEndpoint :one
UPDATE endpoints
SET
    is_active = NOT is_active
WHERE
    id = $1
RETURNING *;

-- name: ListWebhooks :many
SELECT
    id,
    endpoint_id,
    public_key,
    name,
    description,
    is_active,
    created_at,
    updated_at
FROM
    webhooks
WHERE
    endpoint_id = $1
ORDER BY
    updated_at DESC;


-- name: CreateWebhook :one
INSERT INTO webhooks (
    endpoint_id,
    public_key,
    name,
    description
)
VALUES (
    $1, $2, $3, $4
)
RETURNING *;


-- name: ToggleWebhook :one
UPDATE webhooks
SET
    is_active = NOT is_active
WHERE
    id = $1
RETURNING *;
