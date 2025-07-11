package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	AppEnv      string
	DatabaseURL string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env dosyası yüklenemedi")
	}

	cfg := Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL is not set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
