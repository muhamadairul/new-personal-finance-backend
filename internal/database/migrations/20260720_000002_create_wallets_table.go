package migrations

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "20260720_000002_create_wallets_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.Wallet{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&models.Wallet{})
		},
	})
}
