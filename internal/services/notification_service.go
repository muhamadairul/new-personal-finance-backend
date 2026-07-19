package services

import (
	"finance-app-backend/internal/models"
	"finance-app-backend/internal/repositories"
)

type NotificationService struct {
	noteRepo repositories.NotificationRepositoryInterface
}

func NewNotificationService(noteRepo repositories.NotificationRepositoryInterface) *NotificationService {
	return &NotificationService{noteRepo: noteRepo}
}

func (s *NotificationService) List(userID uint, page, limit int) ([]models.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 15
	}
	return s.noteRepo.GetForUser(userID, page, limit)
}

func (s *NotificationService) CountUnread(userID uint) (int64, error) {
	return s.noteRepo.CountUnread(userID)
}

func (s *NotificationService) MarkRead(userID uint, id string) error {
	return s.noteRepo.MarkAsRead(userID, id)
}

func (s *NotificationService) MarkAllRead(userID uint) error {
	return s.noteRepo.MarkAllAsRead(userID)
}
