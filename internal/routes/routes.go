package routes

import (
	"finance-app-backend/internal/controllers"
	"finance-app-backend/internal/database"
	"finance-app-backend/internal/middlewares"
	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/repositories"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

// Setup sets up all the routes for the application
func Setup(app *fiber.App) {
	// Serve uploaded profile photos locally
	app.Static("/storage", "./public/storage")

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, "Application is running smoothly", fiber.Map{
			"status": "UP",
		})
	})

	// Dependency Injection for Auth
	userRepo := repositories.NewUserRepository(database.DB)
	resetRepo := repositories.NewPasswordResetRepository(database.DB)
	authService := services.NewAuthService(userRepo, resetRepo)
	authCtrl := controllers.NewAuthController(authService)

	// API Route Group
	api := app.Group("/api")

	// Public Auth Routes
	api.Post("/register", authCtrl.Register)
	api.Post("/login", authCtrl.Login)
	api.Post("/auth/social", authCtrl.SocialLogin)

	// Forgot Password routes
	api.Post("/password/email", authCtrl.SendOtp)
	api.Post("/password/verify-otp", authCtrl.VerifyOtp)
	api.Post("/password/reset", authCtrl.ResetPassword)

	// Protected Routes (JWT Guarded)
	protected := api.Group("", middlewares.JWTAuth())

	// Profile & Auth
	protected.Post("/logout", authCtrl.Logout)
	protected.Get("/user", authCtrl.GetUser)
	protected.Put("/user/profile", authCtrl.UpdateProfile)
	protected.Post("/user/photo", authCtrl.UploadPhoto)
	protected.Delete("/user/photo", authCtrl.DeletePhoto)
	protected.Post("/user/fcm-token", authCtrl.UpdateFcmToken)
}
