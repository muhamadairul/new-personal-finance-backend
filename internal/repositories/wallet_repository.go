package repositories

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepositoryInterface interface {
	GetForUser(userID uint) ([]models.Wallet, error)
	CountForUser(userID uint) (int64, error)
	GetByID(id uint) (*models.Wallet, error)
	GetByIDWithLock(tx *gorm.DB, id uint) (*models.Wallet, error)
	Create(w *models.Wallet) error
	Update(w *models.Wallet) error
	UpdateWithTx(tx *gorm.DB, w *models.Wallet) error
	Delete(id uint) error
}

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetForUser(userID uint) ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.Where("user_id = ?", userID).Order("created_at ASC").Find(&wallets).Error
	return wallets, err
}

func (r *WalletRepository) CountForUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Wallet{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *WalletRepository) GetByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) GetByIDWithLock(tx *gorm.DB, id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) Create(w *models.Wallet) error {
	return r.db.Create(w).Error
}

func (r *WalletRepository) Update(w *models.Wallet) error {
	return r.db.Save(w).Error
}

func (r *WalletRepository) UpdateWithTx(tx *gorm.DB, w *models.Wallet) error {
	return tx.Save(w).Error
}

func (r *WalletRepository) Delete(id uint) error {
	return r.db.Delete(&models.Wallet{}, id).Error
}
