package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	AppName      string
	AppEnv       string
	AppPort      string
	AppUrl       string
	JwtSecret    string
	AutoMigrate  string
	AutoSeed     string
	DbConnection string
	DbHost       string
	DbPort       string
	DbDatabase   string
	DbUsername   string
	DbPassword   string
	CronSecret   string
}

// AppConfig is the global configuration instance
var AppConfig *Config

// Load loads configurations from .env or environment variables
func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, reading from environment variables")
	}

	AppConfig = &Config{
		AppName:      getEnv("APP_NAME", "Pencatat Keuangan"),
		AppEnv:       getEnv("APP_ENV", "local"),
		AppPort:      getEnv("APP_PORT", "8080"),
		AppUrl:       getEnv("APP_URL", "http://localhost:8080"),
		JwtSecret:    getEnv("JWT_SECRET", "default-secret"),
		AutoMigrate:  getEnv("AUTO_MIGRATE", "true"),
		AutoSeed:     getEnv("AUTO_SEED", "false"),
		DbConnection: getEnv("DB_CONNECTION", "mysql"),
		DbHost:       getEnv("DB_HOST", "127.0.0.1"),
		DbPort:       getEnv("DB_PORT", "3306"),
		DbDatabase:   getEnv("DB_DATABASE", "personal_finance"),
		DbUsername:   getEnv("DB_USERNAME", "root"),
		DbPassword:   getEnv("DB_PASSWORD", ""),
		CronSecret:   getEnv("CRON_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
