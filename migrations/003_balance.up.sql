CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS user_points (
                                           user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    current_balance NUMERIC DEFAULT 0,
    withdrawn_points NUMERIC DEFAULT 0
    );

INSERT INTO user_points (current_balance, withdrawn_points)
VALUES (729.98, 0);
