package middlewares

import (
	"strings"

	"finance-app-backend/internal/config"
	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// JWTAuth guards endpoints using stateless JWT authentication
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Token autentikasi tidak tersedia.", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Format token autentikasi tidak valid.", nil)
		}

		tokenString := parts[1]
		userID, err := utils.VerifyToken(tokenString, config.AppConfig.JwtSecret)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Token autentikasi tidak valid atau telah kadaluarsa.", nil)
		}

		// Store user ID in locals context
		c.Locals("user_id", userID)

		return c.Next()
	}
}
