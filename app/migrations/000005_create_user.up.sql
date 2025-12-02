CREATE TABLE IF NOT EXISTS users (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_type TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at timestamp with time zone not null,

    UNIQUE(email)
);