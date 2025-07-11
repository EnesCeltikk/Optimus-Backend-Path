package models

import "time"

type User struct {
    ID           int       `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    Role         string    `json:"role" db:"role"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Transaction struct {
    ID         int       `json:"id" db:"id"`
    FromUserID int       `json:"from_user_id" db:"from_user_id"`
    ToUserID   int       `json:"to_user_id" db:"to_user_id"`
    Amount     float64   `json:"amount" db:"amount"`
    Type       string    `json:"type" db:"type"`
    Status     string    `json:"status" db:"status"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type Balance struct {
    UserID        int       `json:"user_id" db:"user_id"`
    Amount        float64   `json:"amount" db:"amount"`
    LastUpdatedAt time.Time `json:"last_updated_at" db:"last_updated_at"`
}

type AuditLog struct {
    ID         int             `json:"id" db:"id"`
    EntityType string          `json:"entity_type" db:"entity_type"`
    EntityID   int             `json:"entity_id" db:"entity_id"`
    Action     string          `json:"action" db:"action"`
    Details    map[string]any  `json:"details" db:"details"`
    CreatedAt  time.Time       `json:"created_at" db:"created_at"`
} 