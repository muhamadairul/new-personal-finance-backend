package repositories

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

type TransactionRepositoryInterface interface {
	GetForUser(userID uint, month, year int) ([]models.Transaction, error)
	GetByID(id uint) (*models.Transaction, error)
	Create(tx *gorm.DB, t *models.Transaction) error
	Update(tx *gorm.DB, t *models.Transaction) error
	Delete(tx *gorm.DB, id uint) error
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetForUser(userID uint, month, year int) ([]models.Transaction, error) {
	var txs []models.Transaction
	query := r.db.Where("user_id = ?", userID).
		Preload("Category").
		Preload("Wallet").
		Order("date DESC").
		Order("created_at DESC")

	if month > 0 && year > 0 {
		// Compatible with both MySQL and PostgreSQL extraction
		query = query.Where("EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", month, year)
	}

	err := query.Find(&txs).Error
	return txs, err
}

func (r *TransactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var t models.Transaction
	err := r.db.Preload("Category").Preload("Wallet").First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepository) Create(tx *gorm.DB, t *models.Transaction) error {
	return tx.Create(t).Error
}

func (r *TransactionRepository) Update(tx *gorm.DB, t *models.Transaction) error {
	return tx.Save(t).Error
}

func (r *TransactionRepository) Delete(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.Transaction{}, id).Error
}
