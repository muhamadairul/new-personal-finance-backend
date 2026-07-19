package database

import (
	"fmt"
	"log"

	"finance-app-backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database instance
var DB *gorm.DB

// Connect establishes a connection to the database based on configuration
func Connect() {
	var err error
	cfg := config.AppConfig

	switch cfg.DbConnection {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DbUsername, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbDatabase)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			cfg.DbHost, cfg.DbUsername, cfg.DbPassword, cfg.DbDatabase, cfg.DbPort)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("Database connection type '%s' is not supported", cfg.DbConnection)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Database connection successfully established via %s", cfg.DbConnection)
}
