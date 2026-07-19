package services

import (
	"testing"

	"finance-app-backend/internal/models"
	"finance-app-backend/internal/pkg/fcm"
)

type mockSubRepo struct {
	logs map[uint]*models.SubscriptionLog
	last uint
}

func (m *mockSubRepo) Create(sLog *models.SubscriptionLog) error {
	m.last++
	sLog.ID = m.last
	m.logs[sLog.ID] = sLog
	return nil
}

func (m *mockSubRepo) GetByID(id uint) (*models.SubscriptionLog, error) {
	return m.logs[id], nil
}

func (m *mockSubRepo) GetByExternalID(extID string) (*models.SubscriptionLog, error) {
	for _, l := range m.logs {
		if l.XenditInvoiceID != nil && *l.XenditInvoiceID == extID {
			return l, nil
		}
	}
	return nil, nil
}

func (m *mockSubRepo) GetByXenditInvoiceID(xID string) (*models.SubscriptionLog, error) {
	return nil, nil
}

func (m *mockSubRepo) GetForUser(userID uint) ([]models.SubscriptionLog, error) {
	var res []models.SubscriptionLog
	for _, l := range m.logs {
		if l.UserID == userID {
			res = append(res, *l)
		}
	}
	return res, nil
}

func (m *mockSubRepo) Update(sLog *models.SubscriptionLog) error {
	m.logs[sLog.ID] = sLog
	return nil
}

type mockNoteRepo struct {
	notes map[string]*models.Notification
}

func (m *mockNoteRepo) Create(n *models.Notification) error {
	return nil
}

func (m *mockNoteRepo) GetForUser(userID uint, page, limit int) ([]models.Notification, int64, error) {
	return nil, 0, nil
}

func (m *mockNoteRepo) CountUnread(userID uint) (int64, error) {
	return 0, nil
}

func (m *mockNoteRepo) MarkAsRead(userID uint, id string) error {
	return nil
}

func (m *mockNoteRepo) MarkAllAsRead(userID uint) error {
	return nil
}

func TestSubscriptionService_ActivateSubscription_Stacking(t *testing.T) {
	userRepo := &mockUserRepo{users: make(map[string]*models.User)}
	subRepo := &mockSubRepo{logs: make(map[uint]*models.SubscriptionLog)}
	noteRepo := &mockNoteRepo{notes: make(map[string]*models.Notification)}
	xenditSvc := NewXenditService()
	fcmSvc := fcm.NewFCMService()

	srv := NewSubscriptionService(subRepo, userRepo, noteRepo, xenditSvc, fcmSvc)

	// 1. New User (Expired/Not Pro)
	user := &models.User{
		Name:  "Sub User",
		Email: "sub@example.com",
		IsPro: false,
	}
	userRepo.Create(user)

	planMonthly := "monthly"
	ext1 := "EXT-1"
	log1 := &models.SubscriptionLog{
		UserID:          user.ID,
		Type:            "qris",
		PlanID:          &planMonthly,
		Amount:          15000,
		Status:          "pending",
		XenditInvoiceID: &ext1,
	}
	subRepo.Create(log1)

	// Activate 1st Monthly subscription
	err := srv.ActivateSubscription(log1)
	if err != nil {
		t.Fatalf("ActivateSubscription failed: %v", err)
	}

	if !user.IsPro {
		t.Errorf("expected user to be Pro")
	}

	firstExpiry := *user.SubscriptionUntil

	// 2. User buys 2nd Monthly subscription while 1st is still active -> Stacking
	ext2 := "EXT-2"
	log2 := &models.SubscriptionLog{
		UserID:          user.ID,
		Type:            "qris",
		PlanID:          &planMonthly,
		Amount:          15000,
		Status:          "pending",
		XenditInvoiceID: &ext2,
	}
	subRepo.Create(log2)

	err = srv.ActivateSubscription(log2)
	if err != nil {
		t.Fatalf("ActivateSubscription 2 failed: %v", err)
	}

	expectedSecondExpiry := firstExpiry.AddDate(0, 0, 30)
	if user.SubscriptionUntil.Format("2006-01-02") != expectedSecondExpiry.Format("2006-01-02") {
		t.Errorf("expected stacked expiry %s, got %s", expectedSecondExpiry.Format("2006-01-02"), user.SubscriptionUntil.Format("2006-01-02"))
	}
}
