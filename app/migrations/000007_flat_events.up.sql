CREATE TABLE IF NOT EXISTS outbox (
    id SERIAL PRIMARY KEY,
    event_type TEXT NOT NULL,
    flat_id INTEGER NOT NULL REFERENCES flat(id) ON DELETE CASCADE,
    house_id INTEGER NOT NULL REFERENCES house(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    message TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_outbox_status ON outbox(status, created_at);
