CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE sources (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(), -- public identifier for external use
    egress_url TEXT NOT NULL, -- where deliveries are sent to
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
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (event_id, attempt_number),
    CHECK (
        (state IN ('succeeded', 'failed', 'aborted') AND finished_at IS NOT NULL)
        OR
        (state IN ('pending', 'in_flight') AND finished_at IS NULL)
    )
);
