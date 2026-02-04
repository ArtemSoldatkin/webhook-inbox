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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
