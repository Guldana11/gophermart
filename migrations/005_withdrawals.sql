CREATE TABLE IF NOT EXISTS withdrawals (
                                           id SERIAL PRIMARY KEY,
                                           user_id VARCHAR(36) NOT NULL,
    order_number VARCHAR(255) NOT NULL,
    sum NUMERIC(12,2) NOT NULL,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

DROP TABLE IF EXISTS withdrawals;
