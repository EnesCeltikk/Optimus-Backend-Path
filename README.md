# Optimus Backend Path
Welcome !
---

## 🏗️ Project Structure

- **API Implementation/**: Main Go backend (modular, layered, RESTful)
- **deployment/**: Docker, Compose, Prometheus, Grafana, and monitoring configs
- **db/migrations/**: SQL migration scripts for schema evolution
- **models/**, **services/**, **server/**: Domain-driven, testable, and extendable code

---

## 🚦 Features
- **User, Transaction, Balance domain models**
- **Role-based authentication & authorization** (JWT-ready, dummy token for demo)
- **Concurrent transaction processing** (worker pool, queue, atomic ops)
- **Thread-safe balance updates**
- **Comprehensive API endpoints** (auth, user, transaction, balance)
- **Custom router & middleware** (CORS, rate limit, logging, error handling)
- **Structured logging** (zerolog)
- **Graceful shutdown**
- **Production-ready Docker & Compose setup**
- **Prometheus & Grafana monitoring**

---

## 🚀 Quickstart

### 1. Clone & Configure
```sh
git clone https://github.com/EnesCeltikk/Optimus-Backend-Path.git
cd Optimus-Backend-Path
docker-compose -f deployment/docker-compose.yml up --build
```

- Copy `.env.example` to `.env` and adjust as needed.

### 2. Endpoints
- `POST   /api/v1/auth/register`  — Register user
- `POST   /api/v1/auth/login`     — Login (returns dummy token)
- `GET    /api/v1/users`          — List users (admin)
- `GET    /api/v1/users/{id}`     — Get user by ID
- `PUT    /api/v1/users/{id}`     — Update user
- `DELETE /api/v1/users/{id}`     — Delete user
- `POST   /api/v1/transactions/credit`   — Credit
- `POST   /api/v1/transactions/debit`    — Debit
- `POST   /api/v1/transactions/transfer` — Transfer
- `GET    /api/v1/transactions/history`  — Transaction history
- `GET    /api/v1/transactions/{id}`     — Transaction by ID
- `GET    /api/v1/balances/current`      — Current balance
- `GET    /api/v1/balances/historical`   — Balance history
- `GET    /api/v1/balances/at-time`      — Balance at a specific time

### 3. Monitoring
- **Prometheus**: [http://localhost:9090](http://localhost:9090)
- **Grafana**: [http://localhost:3000](http://localhost:3000) (admin/admin)

---

## 🧩 Architecture
- **Domain-driven**: Models, services, and interfaces are cleanly separated.
- **Concurrent by design**: Transaction queue, worker pool, and thread-safe balance logic.
- **Extensible**: Add new endpoints, services, or middleware with minimal friction.
- **Observability**: Metrics and logs are first-class citizens.

---

## 🛠️ Developer Notes
- **Auth**: Uses a dummy token for demo. Swap in JWT logic for production.
- **Migrations**: Managed via SQL scripts in `db/migrations/`.

---
