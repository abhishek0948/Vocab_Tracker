package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
}

func GetConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "vocab_user"),
		DBPassword: getEnv("DB_PASSWORD", "vocab_password"),
		DBName:     getEnv("DB_NAME", "vocab_tracker"),
		DBPort:     getEnv("DB_PORT", "5432"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}