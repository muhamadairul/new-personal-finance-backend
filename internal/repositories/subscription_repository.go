package repositories

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

type SubscriptionRepositoryInterface interface {
	Create(sLog *models.SubscriptionLog) error
	GetByID(id uint) (*models.SubscriptionLog, error)
	GetByExternalID(extID string) (*models.SubscriptionLog, error)
	GetByXenditInvoiceID(invoiceID string) (*models.SubscriptionLog, error)
	GetForUser(userID uint) ([]models.SubscriptionLog, error)
	Update(sLog *models.SubscriptionLog) error
}

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sLog *models.SubscriptionLog) error {
	return r.db.Create(sLog).Error
}

func (r *SubscriptionRepository) GetByID(id uint) (*models.SubscriptionLog, error) {
	var sLog models.SubscriptionLog
	err := r.db.First(&sLog, id).Error
	if err != nil {
		return nil, err
	}
	return &sLog, nil
}

func (r *SubscriptionRepository) GetByExternalID(extID string) (*models.SubscriptionLog, error) {
	var sLog models.SubscriptionLog
	err := r.db.Where("xendit_invoice_id = ?", extID).First(&sLog).Error
	if err != nil {
		return nil, err
	}
	return &sLog, nil
}

func (r *SubscriptionRepository) GetByXenditInvoiceID(invoiceID string) (*models.SubscriptionLog, error) {
	var sLog models.SubscriptionLog
	err := r.db.Where("xendit_invoice_id = ?", invoiceID).First(&sLog).Error
	if err != nil {
		return nil, err
	}
	return &sLog, nil
}

func (r *SubscriptionRepository) GetForUser(userID uint) ([]models.SubscriptionLog, error) {
	var logs []models.SubscriptionLog
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *SubscriptionRepository) Update(sLog *models.SubscriptionLog) error {
	return r.db.Save(sLog).Error
}
