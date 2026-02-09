CREATE TABLE sources (
    id BIGSERIAL PRIMARY KEY,
    ingress_url TEXT NOT NULL, -- where events arrive from
    egress_url TEXT NOT NULL, -- where deliveries are sent to
    static_headers JSONB NOT NULL DEFAULT '{}',
    status TEXT NOT NULL CHECK (status IN ('active', 'paused', 'quarantined', 'disabled')),
    status_reason VARCHAR(512),
    description VARCHAR(512),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    disable_at TIMESTAMP,
    UNIQUE (ingress_url, egress_url)
);


CREATE TABLE events (

)

CREATE TABLE deliveries (

)






CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(256) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE endpoints (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    name VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    headers JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE webhooks (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id BIGINT REFERENCES endpoints(id) ON DELETE CASCADE,
    public_key TEXT UNIQUE NOT NULL,
    name VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Automatically maintain updated_at on row updates for webhooks
CREATE OR REPLACE FUNCTION set_webhooks_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_webhooks_updated_at
BEFORE UPDATE ON webhooks
FOR EACH ROW
EXECUTE FUNCTION set_webhooks_updated_at();

CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    webhook_id BIGINT REFERENCES webhooks(id) ON DELETE CASCADE,
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    method VARCHAR(16) NOT NULL,
    query_params JSONB,
    headers JSONB,
    body JSONB NOT NULL,
    size INT NOT NULL,
    source_ip INET NOT NULL,
    event_hash TEXT
);

CREATE TABLE deliveries (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT REFERENCES events(id) ON DELETE CASCADE,
    endpoint_id BIGINT REFERENCES endpoints(id) ON DELETE CASCADE,
    status_code INT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
