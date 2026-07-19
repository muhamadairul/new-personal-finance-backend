package seeders

import (
	"log"

	"finance-app-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

func (s *UserSeeder) Run(db *gorm.DB) error {
	var count int64
	adminEmail := "admin@pencatatkeuangan.com"

	db.Model(&models.User{}).Where("email = ?", adminEmail).Count(&count)
	if count > 0 {
		log.Printf("[UserSeeder] Admin user %s already exists, skipping.", adminEmail)
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	if err != nil {
		return err
	}

	pwStr := string(hashedPassword)
	admin := models.User{
		Name:     "Admin",
		Email:    adminEmail,
		Password: &pwStr,
		IsAdmin:  true,
	}

	return db.Create(&admin).Error
}

func init() {
	Register(&UserSeeder{})
}
