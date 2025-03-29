CREATE TABLE payments (
    idempotency_key TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    amount BIGINT NOT NULL,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('in_progress', 'success')),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,             
    event_type TEXT NOT NULL,          
    payload JSONB NOT NULL,            
    status TEXT NOT NULL CHECK(status IN ('new', 'success', 'in_progress')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);