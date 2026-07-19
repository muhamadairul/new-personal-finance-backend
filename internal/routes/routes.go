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

	// Dependency Injection for Category
	categoryRepo := repositories.NewCategoryRepository(database.DB)
	categoryService := services.NewCategoryService(categoryRepo, userRepo)
	categoryCtrl := controllers.NewCategoryController(categoryService)

	// Dependency Injection for Wallet
	walletRepo := repositories.NewWalletRepository(database.DB)
	walletService := services.NewWalletService(walletRepo, userRepo)
	walletCtrl := controllers.NewWalletController(walletService)

	// Dependency Injection for Transaction
	txRepo := repositories.NewTransactionRepository(database.DB)
	txService := services.NewTransactionService(txRepo, walletRepo, categoryRepo, database.DB)
	txCtrl := controllers.NewTransactionController(txService)

	// Dependency Injection for Budget
	budgetRepo := repositories.NewBudgetRepository(database.DB)
	budgetService := services.NewBudgetService(budgetRepo, categoryRepo, database.DB)
	budgetCtrl := controllers.NewBudgetController(budgetService)

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

	// Categories CRUD
	protected.Get("/categories", categoryCtrl.Index)
	protected.Post("/categories", categoryCtrl.Store)
	protected.Get("/categories/:id", categoryCtrl.Show)
	protected.Put("/categories/:id", categoryCtrl.Update)
	protected.Delete("/categories/:id", categoryCtrl.Destroy)

	// Wallets CRUD
	protected.Get("/wallets", walletCtrl.Index)
	protected.Post("/wallets", walletCtrl.Store)
	protected.Get("/wallets/:id", walletCtrl.Show)
	protected.Put("/wallets/:id", walletCtrl.Update)
	protected.Delete("/wallets/:id", walletCtrl.Destroy)

	// Transactions CRUD
	protected.Get("/transactions", txCtrl.Index)
	protected.Post("/transactions", txCtrl.Store)
	protected.Get("/transactions/:id", txCtrl.Show)
	protected.Put("/transactions/:id", txCtrl.Update)
	protected.Delete("/transactions/:id", txCtrl.Destroy)

	// Budgets CRUD
	protected.Get("/budgets", budgetCtrl.Index)
	protected.Post("/budgets", budgetCtrl.Store)
	protected.Get("/budgets/:id", budgetCtrl.Show)
	protected.Put("/budgets/:id", budgetCtrl.Update)
	protected.Delete("/budgets/:id", budgetCtrl.Destroy)
}
