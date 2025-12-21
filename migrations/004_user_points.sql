CREATE TABLE IF NOT EXISTS user_points (
                                           user_id UUID PRIMARY KEY,
                                           current_balance NUMERIC(12, 2) NOT NULL DEFAULT 0,
    withdrawn_points NUMERIC(12, 2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user_points_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
                         ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_user_points_user_id
    ON user_points(user_id);
