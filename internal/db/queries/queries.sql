-- name: ListSources :many
SELECT
    id
    , public_id
    , egress_url
    , static_headers
    , status
    , status_reason
    , description
    , created_at
    , updated_at
    , disable_at
FROM
    sources
WHERE
    (
        @cursor_ts::timestamptz IS NULL OR
        (
            updated_at < @cursor_ts::timestamptz OR
            (
                updated_at = @cursor_ts::timestamptz AND
                id < @cursor_id
            )
        )
    ) AND
    (
        @search_query::text IS NULL OR
        @search_query::text = '' OR
        (
            egress_url ILIKE '%' || @search_query::text || '%' OR
            description ILIKE '%' || @search_query::text || '%' OR
            public_id::text ILIKE '%' || @search_query::text || '%'
        )
    ) AND
    (
        @filter_status::text IS NULL OR
        @filter_status::text = '' OR
        @filter_status::text = '*' OR
        status = @filter_status::text
    )
ORDER BY
    updated_at DESC
    , id DESC
LIMIT
    @page_size + 1;

-- name: GetSourceByID :one
SELECT
    id
    , public_id
    , egress_url
    , static_headers
    , status
    , status_reason
    , description
    , created_at
    , updated_at
    , disable_at
FROM
    sources
WHERE
    id = @source_id;

-- name: GetSourceByPublicID :one
SELECT
    id
    , public_id
    , egress_url
    , static_headers
    , status
    , status_reason
    , description
    , created_at
    , updated_at
    , disable_at
FROM
    sources
WHERE
    public_id = @public_id;

-- name: CreateSource :one
INSERT INTO sources (
    egress_url
    , static_headers
    , description
) VALUES (
    @egress_url
    , @static_headers
    , @description
)
RETURNING *;

-- name: ListEventsBySource :many
SELECT
    id
    , source_id
    , dedup_hash
    , method
    , ingress_path
    , remote_address
    , query_params
    , raw_headers
    , body
    , body_content_type
    , received_at
FROM
    events
WHERE
    (
        source_id = @source_id AND
        (
            @cursor_ts::timestamptz IS NULL OR
            received_at < @cursor_ts::timestamptz OR
            (
                received_at = @cursor_ts::timestamptz AND
                id < @cursor_id
            )
        )
    ) AND
    (
        @search_query::text IS NULL OR @search_query::text = '' OR
        (
            dedup_hash ILIKE '%' || @search_query::text || '%' OR
            method ILIKE '%' || @search_query::text || '%' OR
            ingress_path ILIKE '%' || @search_query::text || '%' OR
            remote_address::text ILIKE '%' || @search_query::text || '%'
        )
    )
ORDER BY
    received_at DESC
    , id DESC
LIMIT
    @page_size + 1;


-- name: GetEventByID :one
SELECT
    id
    , source_id
    , dedup_hash
    , method
    , ingress_path
    , remote_address
    , query_params
    , raw_headers
    , body
    , body_content_type
    , received_at
FROM
    events
WHERE
    id = @event_id;


-- name: CreateEvent :one
INSERT INTO events (
    source_id
    , dedup_hash
    , method
    , ingress_path
    , remote_address
    , query_params
    , raw_headers
    , body
    , body_content_type
) VALUES (
    @source_id
    , @dedup_hash
    , @method
    , @ingress_path
    , @remote_address
    , @query_params
    , @raw_headers
    , @body
    , @body_content_type
)
RETURNING id;

-- name: ListDeliveryAttemptsByEvent :many
SELECT
    id
    , event_id
    , attempt_number
    , state
    , status_code
    , error_type
    , error_message
    , started_at
    , finished_at
    , created_at
    , next_attempt_at
FROM
    delivery_attempts
WHERE
    (
        event_id = @event_id AND
        (
            @cursor_ts::timestamptz IS NULL OR
            created_at < @cursor_ts::timestamptz OR
            (
                created_at = @cursor_ts::timestamptz AND
                id < @cursor_id
            )
        )
    ) AND
    (
        @search_query::text IS NULL OR
        @search_query::text = '' OR
        (
            state ILIKE '%' || @search_query::text || '%' OR
            status_code::text ILIKE '%' || @search_query::text || '%' OR
            error_type ILIKE '%' || @search_query::text || '%' OR
            error_message ILIKE '%' || @search_query::text || '%'
        )
    ) AND
    (
        @filter_state::text IS NULL OR
        @filter_state::text = '' OR
        @filter_state::text = '*' OR
        state = @filter_state::text
    )
ORDER BY
    created_at DESC
    , id DESC
LIMIT
    @page_size + 1;

-- name: CreateDeliveryAttempt :one
INSERT INTO delivery_attempts (
    event_id
    , attempt_number
    , state
    , status_code
    , error_type
    , error_message
    , started_at
    , finished_at
    , next_attempt_at
) VALUES (
    @event_id
    , @attempt_number
    , @state
    , @status_code
    , @error_type
    , @error_message
    , @started_at
    , @finished_at
    , @next_attempt_at
)
RETURNING id;


-- name: SelectPendingDeliveryAttemptIDs :many
SELECT
    delivery_attempts.id
FROM
    delivery_attempts
INNER JOIN events
    ON
        delivery_attempts.event_id = events.id
INNER JOIN sources
    ON
        events.source_id = sources.id
WHERE
    delivery_attempts.event_id = events.id AND
    delivery_attempts.state = 'pending' AND
    sources.status = 'active' AND
    COALESCE(delivery_attempts.next_attempt_at, NOW()) <= NOW()
ORDER BY
    delivery_attempts.created_at ASC
FOR UPDATE OF delivery_attempts SKIP LOCKED
LIMIT
    @batch_size;


-- name: UpdateDeliveryAttemptsToInFlight :many
UPDATE delivery_attempts
SET
    state = 'in_flight'
    , started_at = NOW()
    , finished_at = NULL
WHERE
    id = ANY(@batch_ids::bigint[])
RETURNING
    id
    , event_id
    , attempt_number;

-- name: UpdateDeliveryAttempt :exec
UPDATE delivery_attempts
SET
    state = @state
    , status_code = @status_code
    , error_type = @error_type
    , error_message = @error_message
    , started_at = CASE
        WHEN @state = 'pending' THEN NULL
        ELSE COALESCE(started_at, @started_at)
    END
    , finished_at = CASE
        WHEN @state IN ('pending', 'in_flight') THEN NULL
        ELSE COALESCE(finished_at, @finished_at)
    END
WHERE
    id = @delivery_attempt_id AND
    state IN ('pending', 'in_flight');

-- name: RecoverStuckDeliveryAttempts :exec
UPDATE delivery_attempts
SET
    state = 'pending'
    , started_at = NULL
    , finished_at = NULL
WHERE
    state = 'in_flight' AND
    started_at < NOW() - INTERVAL '15 minutes';
