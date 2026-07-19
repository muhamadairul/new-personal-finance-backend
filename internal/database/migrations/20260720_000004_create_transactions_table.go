package migrations

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "20260720_000004_create_transactions_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.Transaction{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&models.Transaction{})
		},
	})
}
