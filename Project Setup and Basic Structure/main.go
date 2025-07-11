package main

import (
	"log"

	"myapp/config"
	"myapp/db"
	"myapp/server"
)

func main() {
	cfg := config.Load()

	db.Migrate(cfg.DatabaseURL, "db/migrations")

	database := db.Connect(cfg.DatabaseURL)

	log.Println("Uygulama başladı...")

	server.StartServer(cfg.AppPort)

	_ = database
}
