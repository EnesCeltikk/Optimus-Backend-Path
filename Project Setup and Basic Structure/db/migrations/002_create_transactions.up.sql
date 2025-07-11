CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_user_id INT REFERENCES users(id),
    to_user_id INT REFERENCES users(id),
    amount NUMERIC(12, 2) NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_from_user ON transactions(from_user_id);
CREATE INDEX idx_transactions_to_user ON transactions(to_user_id);
