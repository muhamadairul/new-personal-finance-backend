package migrations

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "20260720_000008_create_notifications_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.Notification{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&models.Notification{})
		},
	})
}
