package migrations

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "20260720_000007_create_password_reset_tokens_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.PasswordResetToken{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&models.PasswordResetToken{})
		},
	})
}
