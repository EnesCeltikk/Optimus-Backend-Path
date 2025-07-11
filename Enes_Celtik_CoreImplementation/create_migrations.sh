#!/usr/bin/env bash

mkdir -p db/migrations

# users tablosu için migration dosyalarını oluştur
cat <<EOF > db/migrations/001_create_users.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
EOF

cat <<EOF > db/migrations/001_create_users.down.sql
DROP TABLE IF EXISTS users;
EOF

# transactions tablosu için migration dosyalarını oluştur
cat <<EOF > db/migrations/002_create_transactions.up.sql
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
EOF

cat <<EOF > db/migrations/002_create_transactions.down.sql
DROP TABLE IF EXISTS transactions;
EOF

# balances tablosu için migration dosyalarını oluştur
cat <<EOF > db/migrations/003_create_balances.up.sql
CREATE TABLE balances (
    user_id INT PRIMARY KEY REFERENCES users(id),
    amount NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
EOF

cat <<EOF > db/migrations/003_create_balances.down.sql
DROP TABLE IF EXISTS balances;
EOF

# audit_logs tablosu için migration dosyalarını oluştur
cat <<EOF > db/migrations/004_create_audit_logs.up.sql
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT NOT NULL,
    action VARCHAR(50) NOT NULL,
    details TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
EOF

cat <<EOF > db/migrations/004_create_audit_logs.down.sql
DROP TABLE IF EXISTS audit_logs;
EOF

echo "Migration dosyaları oluşturuldu!"
