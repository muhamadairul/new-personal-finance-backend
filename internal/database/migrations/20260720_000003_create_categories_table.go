package migrations

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "20260720_000003_create_categories_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.Category{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&models.Category{})
		},
	})
}
