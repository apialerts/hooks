-- +goose Up

CREATE TABLE endpoints (
    id TEXT PRIMARY KEY,
    creator_ip TEXT NOT NULL,
    response_status INTEGER DEFAULT 200,
    response_delay_ms INTEGER DEFAULT 0,
    response_body TEXT DEFAULT '',
    request_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    last_activity_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_endpoints_last_activity ON endpoints(last_activity_at);
CREATE INDEX idx_endpoints_creator_ip ON endpoints(creator_ip);

CREATE TABLE requests (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id TEXT NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
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

CREATE INDEX idx_requests_endpoint_id ON requests(endpoint_id, received_at DESC);

-- +goose Down
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS endpoints;
