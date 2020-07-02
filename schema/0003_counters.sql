CREATE TABLE counters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    value BIGINT NOT NULL
);
