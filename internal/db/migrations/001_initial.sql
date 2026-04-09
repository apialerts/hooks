-- +goose Up

CREATE TABLE endpoint (
    id BIGSERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    creator_ip TEXT NOT NULL,
    response_status INTEGER DEFAULT 200,
    response_delay_ms INTEGER DEFAULT 0,
    request_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    last_activity_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_endpoint_last_activity ON endpoint(last_activity_at);
CREATE INDEX idx_endpoint_creator_ip ON endpoint(creator_ip);

CREATE TABLE request (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id BIGINT NOT NULL REFERENCES endpoint(id) ON DELETE CASCADE,
    seq INTEGER DEFAULT 0,
    method TEXT NOT NULL,
    path TEXT,
    query_params TEXT,
    headers JSONB NOT NULL,
    body TEXT,
    body_size INTEGER NOT NULL,
    content_type TEXT,
    source_ip TEXT,
    response_status INTEGER,
    response_body TEXT DEFAULT '',
    received_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_request_endpoint_id ON request(endpoint_id, received_at DESC);

CREATE TABLE endpoint_response (
    endpoint_id BIGINT NOT NULL REFERENCES endpoint(id) ON DELETE CASCADE,
    status INTEGER NOT NULL,
    delay_ms INTEGER NOT NULL DEFAULT 0,
    body TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (endpoint_id, status, delay_ms)
);

-- +goose Down
DROP TABLE IF EXISTS endpoint_response;
DROP TABLE IF EXISTS request;
DROP TABLE IF EXISTS endpoint;
