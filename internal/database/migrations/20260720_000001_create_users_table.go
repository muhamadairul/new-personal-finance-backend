package migrations

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "20260720_000001_create_users_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.User{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&models.User{})
		},
	})
}
