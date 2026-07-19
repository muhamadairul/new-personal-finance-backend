package response

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents uniform response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Success returns success response. Always uses HTTP 200 OK.
func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessNoContent returns success response without data. Always uses HTTP 200 OK.
func SuccessNoContent(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
	})
}

// Error returns error response with custom HTTP status code.
func Error(c *fiber.Ctx, status int, message string, errors interface{}) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
