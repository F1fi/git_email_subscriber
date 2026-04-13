CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    repo TEXT NOT NULL,
    confirmed BOOLEAN DEFAULT false,
    confirm_token TEXT UNIQUE,
    unsubscribe_token TEXT UNIQUE,
    last_seen_tag TEXT,
    created_at TIMESTAMP DEFAULT now()
);