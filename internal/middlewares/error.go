package middlewares

import (
	"log"

	"finance-app-backend/internal/config"
	"finance-app-backend/internal/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// GlobalErrorHandler catches all panic/errors in request handlers
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// If it's a Fiber error, retrieve the status code
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Printf("[Error Handler] Route: %s, Error: %v", c.Path(), err)

	// In local dev, reveal detailed messages. In production, hide it.
	message := "Terjadi kesalahan internal pada server."
	if config.AppConfig.AppEnv == "local" {
		message = err.Error()
	}

	return response.Error(c, code, message, nil)
}
