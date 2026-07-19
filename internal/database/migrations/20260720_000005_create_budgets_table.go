package migrations

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "20260720_000005_create_budgets_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.Budget{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&models.Budget{})
		},
	})
}
