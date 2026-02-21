CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE sources (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(), -- public identifier for external use
    egress_url TEXT NOT NULL CHECK (
        -- Require HTTP/HTTPS and block obvious internal/loopback/metadata targets for SSRF mitigation
        egress_url ~ '^https?://' AND
        egress_url !~* '^https?://(localhost|127\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})|0\.0\.0\.0|\[?::1\]?)(/|:|$)' AND
        egress_url !~* '^https?://\[\:\:ffff\:127\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\]' AND
        egress_url !~* '^https?://10\.' AND
        egress_url !~* '^https?://192\.168\.' AND
        egress_url !~* '^https?://172\.(1[6-9]|2[0-9]|3[0-1])\.' AND
        egress_url !~* '^https?://169\.254\.169\.254(/|:|$)' AND
        egress_url !~* '^https?://\[::ffff:0\.0\.0\.0\]' AND
        egress_url !~* '^https?://localhost\.(/|:|$)' AND
        CHAR_LENGTH(egress_url) <= 2048
    ), -- where deliveries are sent to
    static_headers JSONB NOT NULL DEFAULT '{}',
    status TEXT NOT NULL CHECK (status IN ('active', 'paused', 'quarantined', 'disabled')) DEFAULT 'active',
    status_reason VARCHAR(512),
    description VARCHAR(512),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    disable_at TIMESTAMPTZ
);

CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    source_id BIGINT NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    dedup_hash TEXT,
    method VARCHAR(16) NOT NULL,
    ingress_path TEXT NOT NULL,
    remote_address INET,
    query_params JSONB NOT NULL DEFAULT '{}',
    raw_headers JSONB NOT NULL DEFAULT '{}',
    body BYTEA NOT NULL,
    body_content_type TEXT NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE delivery_attempts (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL
        REFERENCES events(id) ON DELETE CASCADE,
    attempt_number INT NOT NULL,
    state TEXT NOT NULL
        CHECK (state IN ('pending', 'in_flight', 'succeeded', 'failed', 'aborted')),
    status_code INT,
    error_type TEXT,
    error_message TEXT,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    next_attempt_at TIMESTAMPTZ,
    UNIQUE (event_id, attempt_number),
    CHECK (
        (state IN ('succeeded', 'failed', 'aborted') AND finished_at IS NOT NULL)
        OR
        (state IN ('pending', 'in_flight') AND finished_at IS NULL)
    )
);
