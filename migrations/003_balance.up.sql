CREATE TABLE IF NOT EXISTS user_points (
                                           user_id UUID PRIMARY KEY,
                                           current_balance NUMERIC DEFAULT 0,
                                           withdrawn_points NUMERIC DEFAULT 0
);
