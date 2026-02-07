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

-- name: GetEndpointByID :one
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
    id = $1;


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
    is_active = NOT is_active,
    updated_at = NOW()
WHERE
    id = $1
RETURNING *;


-- name: ListEvents :many
SELECT
    id,
    webhook_id,
    received_at,
    method,
    query_params,
    headers,
    body,
    size,
    source_ip,
    event_hash
FROM
    events
WHERE
    webhook_id = $1
ORDER BY
    received_at DESC;

-- name: ListDeliveries :many
SELECT
    id,
    event_id,
    endpoint_id,
    status_code,
    error_message,
    created_at
FROM
    deliveries
WHERE
    endpoint_id = $1
ORDER BY
    created_at DESC;
