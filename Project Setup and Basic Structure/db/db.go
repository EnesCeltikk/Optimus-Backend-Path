package db

import (
	"database/sql"
	"fmt"
	"log"

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


func Connect(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}
	fmt.Println("✅ Connected to database")
	return db
}
