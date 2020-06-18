CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    content TEXT NOT NULL,
    external_id UUID
);
