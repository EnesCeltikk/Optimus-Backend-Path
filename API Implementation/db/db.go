package db

import (
	"database/sql"
	"fmt"
	"log"

	"CoreImplementation/models"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate(dbURL, migrationsPath string) {
    m, err := migrate.New(
        "file://"+migrationsPath,
        dbURL,
    )
    if err != nil {
        log.Fatalf("Migration init error: %v", err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Migration up error: %v", err)
    }
    fmt.Println("✅ Migrations applied")
}

func Connect(dbURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, fmt.Errorf("DB connection error: %v", err)
    }
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("DB ping error: %v", err)
    }
    fmt.Println("✅ Connected to database")
    return db, nil
}

func InsertUser(db *sql.DB, user *models.User) error {
    _, err := db.Exec(
        "INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4)",
        user.Username, user.Email, user.PasswordHash, user.Role,
    )
    return err
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
    row := db.QueryRow(
        "SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE email=$1",
        email,
    )

    var u models.User
    if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
        return nil, fmt.Errorf("failed to get user: %v", err)
    }
    return &u, nil
}

func InsertTransaction(db *sql.DB, tx *models.Transaction) error {
    _, err := db.Exec(
        "INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
        tx.FromUserID, tx.ToUserID, tx.Amount, tx.Type, tx.Status, tx.CreatedAt,
    )
    return err
}

func UpdateTransactionStatus(db *sql.DB, id int, status string) error {
    _, err := db.Exec(
        "UPDATE transactions SET status = $1 WHERE id = $2",
        status, id,
    )
    return err
}

func GetBalance(db *sql.DB, userID int) (*models.Balance, error) {
    row := db.QueryRow(
        "SELECT user_id, amount, last_updated_at FROM balances WHERE user_id = $1",
        userID,
    )
    var b models.Balance
    if err := row.Scan(&b.UserID, &b.Amount, &b.LastUpdatedAt); err != nil {
        return nil, fmt.Errorf("failed to get balance: %v", err)
    }
    return &b, nil
}

func UpdateBalance(db *sql.DB, b *models.Balance) error {
    _, err := db.Exec(
        "UPDATE balances SET amount = $1, last_updated_at = $2 WHERE user_id = $3",
        b.Amount, b.LastUpdatedAt, b.UserID,
    )
    return err
}
