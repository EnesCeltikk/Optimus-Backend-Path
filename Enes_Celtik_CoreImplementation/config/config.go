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
	// .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env dosyası yüklenemedi")
	}

	// Yapılandırmaları al
	cfg := Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	// DATABASE_URL boşsa hata ver ve uygulamayı sonlandır
	if cfg.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL is not set in the .env file")
	}

	return cfg
}

// Environment variable'larını kontrol eder ve değer döner
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
