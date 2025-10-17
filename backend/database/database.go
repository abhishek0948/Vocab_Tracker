package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/abhishek0948/Vocab_Tracker/config"
	"github.com/abhishek0948/Vocab_Tracker/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ensureSSLMode(dsn string) string {
	if strings.Contains(dsn, "sslmode=") {
		return dsn
	}

	if strings.Contains(dsn, "?") {
		return dsn + "&sslmode=require"
	}
	return dsn + "?sslmode=require"
}

func Connect() {
	cfg := config.GetConfig()

	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		dsn = ensureSSLMode(dsn)
		log.Println("Using DATABASE_URL from environment for DB connection")
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
		log.Println("Using individual DB config fields for connection")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql DB from gorm DB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Database connected and reachable")
}

func Migrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Vocabulary{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")
}