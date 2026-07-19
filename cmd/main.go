package main

import (
	"log"

	"finance-app-backend/internal/config"
	"finance-app-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load config
	config.Load()

	// Connect to database
	// NOTE: Uncomment database.Connect() once DB credentials are set up in .env
	// database.Connect()

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
