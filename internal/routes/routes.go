package routes

import (
	"finance-app-backend/internal/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// Setup sets up all the routes for the application
func Setup(app *fiber.App) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, "Application is running smoothly", fiber.Map{
			"status": "UP",
		})
	})
}
