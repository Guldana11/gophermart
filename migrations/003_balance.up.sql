CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS user_points (
                                           user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    current_balance NUMERIC DEFAULT 0,
    withdrawn_points NUMERIC DEFAULT 0
    );

INSERT INTO user_points (user_id, current_balance, withdrawn_points)
VALUES ('11111111-1111-1111-1111-111111111111', 729.98, 0)
    ON CONFLICT (user_id) DO UPDATE
                                 SET current_balance = EXCLUDED.current_balance,
                                 withdrawn_points = EXCLUDED.withdrawn_points;
