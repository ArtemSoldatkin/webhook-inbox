-- name: ListSources :many
SELECT
    id,
    public_id,
    egress_url,
    static_headers,
    status,
    status_reason,
    description,
    created_at,
    updated_at,
    disable_at
FROM
    sources
ORDER BY
    created_at DESC;

-- name: CreateSource :one
INSERT INTO sources (
    egress_url,
    static_headers,
    description
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: ListEventsBySource :many
SELECT
    id,
    source_id,
    dedup_hash,
    method,
    query_params,
    raw_headers,
    body,
    body_content_type,
    received_at
FROM
    events
WHERE
    source_id = $1
ORDER BY
    received_at DESC;


-- name: ListDeliveryAttemptsByEvent :many
SELECT
    id,
    event_id,
    attempt_number,
    state,
    status_code,
    error_type,
    error_message,
    started_at,
    finished_at,
    created_at
FROM
    delivery_attempts
WHERE
    event_id = $1
ORDER BY
    created_at DESC;
