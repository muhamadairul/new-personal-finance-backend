package migrations

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "20260720_000009_create_api_logs_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.ApiLog{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&models.ApiLog{})
		},
	})
}
