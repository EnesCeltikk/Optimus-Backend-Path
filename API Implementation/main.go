package main

import (
	"CoreImplementation/config"
	"CoreImplementation/db"
	"CoreImplementation/server"
	"CoreImplementation/services"
	"log"
)

func main() {
	cfg := config.Load()

	db.Migrate(cfg.DatabaseURL, "db/migrations")

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Veritabanı bağlantı hatası: %v", err)
	}

	userService := services.NewUserService(database)
	transactionService := services.NewTransactionService(database)
	balanceService := services.NewBalanceService(database)

	username := "john_doe"
	email := "john@example.com"
	password := "securepassword"

	log.Println("== Kullanıcı Kaydı Testi ==")
	user, err := userService.Register(username, email, password)
	if err != nil {
		if err.Error() == "user already exists" {
			log.Printf("Uyarı: Kullanıcı zaten mevcut. E-posta: %s", email)
		} else {
			log.Fatalf("Kullanıcı kaydı hatası: %v", err)
		}
	} else {
		log.Printf("Yeni kullanıcı kaydedildi: %+v\n", user)
	}

	log.Println("== Tüm Kullanıcıları Listeleme ==")
	users, err := userService.GetAllUsers()
	if err != nil {
		log.Fatalf("Kullanıcı listeleme hatası: %v", err)
	}
	for _, u := range users {
		log.Printf("Kullanıcı: ID=%d, Username=%s, Email=%s\n", u.ID, u.Username, u.Email)
	}

	log.Println("== Kullanıcı Silme Testi ==")
	err = userService.DeleteUserByEmail(email)
	if err != nil {
		log.Fatalf("Kullanıcı silme hatası: %v", err)
	} else {
		log.Printf("Kullanıcı silindi: %s\n", email)
	}

	server.StartServer(cfg.AppPort, userService, transactionService, balanceService)
}
