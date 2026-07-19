package main

import (
	"log"

	"finance-app-backend/internal/config"
	"finance-app-backend/internal/database"
	"finance-app-backend/internal/database/migrations"
	"finance-app-backend/internal/database/seeders"
	"finance-app-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load config
	config.Load()

	// Connect to database if configured
	if config.AppConfig.DbHost != "" {
		database.Connect()

		if config.AppConfig.AutoMigrate == "true" {
			log.Println("AUTO_MIGRATE=true: Running database migrations...")
			if err := migrations.Migrate(database.DB); err != nil {
				log.Printf("Auto migration error: %v", err)
			}
		}

		if config.AppConfig.AutoSeed == "true" {
			log.Println("AUTO_SEED=true: Running database seeders...")
			if err := seeders.RunAll(database.DB, ""); err != nil {
				log.Printf("Auto seed error: %v", err)
			}
		}
	}

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: config.AppConfig.AppName,
	})

	// Register Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Register Routes
	routes.Setup(app)

	// Start server
	log.Printf("Server starting on port %s...", config.AppConfig.AppPort)
	if err := app.Listen(":" + config.AppConfig.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
