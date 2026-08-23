package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	DatabaseURL         string
	FrontendURL         string
	JWTSecret           string
	JWTExpiryHours      string
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	StripeSecretKey     string
	StripeWebhookSecret string
}

// Load reads .env file and returns a Config struct.
// Calling os.Getenv directly after godotenv.Load ensures
// all values come from one source of truth.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	return &Config{
		Port:                getEnv("PORT", "8080"),
		Env:                 getEnv("ENV", "development"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		FrontendURL:         getEnv("FRONTEND_URL", ""),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		JWTExpiryHours:      getEnv("JWT_EXPIRY_HOURS", "24"),
		GoogleClientID:      getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:  getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:   getEnv("GOOGLE_REDIRECT_URL", ""),
		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
	}
}

// getEnv returns the value of an environment variable or a fallback default.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
