CREATE TABLE IF NOT EXISTS withdrawals (
                                           order_number VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_points(user_id) ON DELETE CASCADE,
    sum NUMERIC(12,2) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

