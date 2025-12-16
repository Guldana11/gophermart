DROP TABLE IF EXISTS orders;

CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        number TEXT UNIQUE NOT NULL,
                        user_id TEXT NOT NULL,
                        status TEXT NOT NULL DEFAULT 'NEW',
                        accrual NUMERIC(10,2) DEFAULT 0,
                        uploaded_at TIMESTAMP NOT NULL DEFAULT now(),
                        CONSTRAINT chk_number_digits CHECK (number ~ '^[0-9]+$')
    );
