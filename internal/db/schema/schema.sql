CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(256) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE endpoints (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    public_key TEXT NOT NULL,
    name VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_received_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_endpoints_user_id_public_key ON endpoints(user_id, public_key);

CREATE INDEX idx_endpoints_user_id ON endpoints(user_id);

CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id BIGINT REFERENCES endpoints(id) ON DELETE CASCADE,
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    headers JSONB,
    body JSONB NOT NULL,
    size INT NOT NULL,
    source_ip INET NOT NULL,
    content_type TEXT NOT NULL,
    event_hash TEXT
);

CREATE INDEX idx_events_endpoint_id_received_at_id ON events(endpoint_id, received_at DESC, id DESC);
CREATE INDEX idx_events_event_hash ON events(event_hash);
