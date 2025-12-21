CREATE TABLE IF NOT EXISTS withdrawals (
    order_id TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    sum NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
