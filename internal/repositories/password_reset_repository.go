package repositories

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

// PasswordResetRepositoryInterface specifies password tokens DB operations
type PasswordResetRepositoryInterface interface {
	Create(token *models.PasswordResetToken) error
	GetByEmail(email string) (*models.PasswordResetToken, error)
	DeleteByEmail(email string) error
}

type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *PasswordResetRepository) GetByEmail(email string) (*models.PasswordResetToken, error) {
	var record models.PasswordResetToken
	err := r.db.Where("email = ?", email).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *PasswordResetRepository) DeleteByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.PasswordResetToken{}).Error
}
