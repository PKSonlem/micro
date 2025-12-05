CREATE TABLE IF NOT EXISTS subscription (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    house_id INTEGER NOT NULL REFERENCES house(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(email, house_id)
);

CREATE INDEX idx_subscription_house_id ON subscription(house_id);