package repositories

import (
	"time"

	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

type NotificationRepositoryInterface interface {
	Create(n *models.Notification) error
	GetForUser(userID uint, page, limit int) ([]models.Notification, int64, error)
	CountUnread(userID uint) (int64, error)
	MarkAsRead(userID uint, id string) error
	MarkAllAsRead(userID uint) error
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *models.Notification) error {
	return r.db.Create(n).Error
}

func (r *NotificationRepository) GetForUser(userID uint, page, limit int) ([]models.Notification, int64, error) {
	var notes []models.Notification
	var total int64

	r.db.Model(&models.Notification{}).Where("notifiable_id = ?", userID).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Where("notifiable_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notes).Error

	return notes, total, err
}

func (r *NotificationRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("notifiable_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *NotificationRepository) MarkAsRead(userID uint, id string) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("id = ? AND notifiable_id = ?", id, userID).
		Update("read_at", &now).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("notifiable_id = ? AND read_at IS NULL", userID).
		Update("read_at", &now).Error
}
