package migrations

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "20260720_000006_create_subscription_logs_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.SubscriptionLog{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&models.SubscriptionLog{})
		},
	})
}
