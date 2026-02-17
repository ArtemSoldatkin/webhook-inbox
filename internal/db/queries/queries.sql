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

-- name: GetSourceByID :one
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
WHERE
    id = $1;

-- name: GetSourceByPublicID :one
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
WHERE
    public_id = $1;

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
    ingress_path,
    remote_address,
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


-- name: GetEventByID :one
SELECT
    id,
    source_id,
    dedup_hash,
    method,
    ingress_path,
    remote_address,
    query_params,
    raw_headers,
    body,
    body_content_type,
    received_at
FROM
    events
WHERE
    id = $1;


-- name: CreateEvent :one
INSERT INTO events (
    source_id,
    dedup_hash,
    method,
    ingress_path,
    remote_address,
    query_params,
    raw_headers,
    body,
    body_content_type
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING id;

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

-- name: CreateDeliveryAttempt :one
INSERT INTO delivery_attempts (
    event_id,
    attempt_number,
    state,
    status_code,
    error_type,
    error_message,
    started_at,
    finished_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING id;

-- name: ListPendingDeliveryAttempts :many
SELECT
    delivery_attempts.id,
    delivery_attempts.event_id,
    delivery_attempts.attempt_number
FROM
    delivery_attempts
INNER JOIN events
    ON delivery_attempts.event_id = events.id
INNER JOIN sources
    ON events.source_id = sources.id
WHERE
    delivery_attempts.state = 'pending'
    AND sources.status = 'active'
ORDER BY
    delivery_attempts.created_at ASC
FOR UPDATE SKIP LOCKED;


-- name: UpdateDeliveryAttempt :exec
UPDATE delivery_attempts
SET
    state = $1,
    status_code = $2,
    error_type = $3,
    error_message = $4,
    started_at = CASE
        WHEN $1 = 'pending' THEN NULL
        ELSE COALESCE(started_at, $5)
    END,
    finished_at = CASE
        WHEN $1 IN ('pending', 'in_flight') THEN NULL
        ELSE COALESCE(finished_at, $6)
    END
WHERE
    id = $7
    AND state IN ('pending', 'in_flight');

-- name: RecoverStuckDeliveryAttempts :exec
UPDATE delivery_attempts
SET
    state = 'pending',
    started_at = NULL,
    finished_at = NULL
WHERE
    state = 'in_flight'
    AND started_at < NOW() - INTERVAL '15 minutes';
