package services

import (
	"errors"
	"testing"

	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

type mockWalletRepo struct {
	wallets map[uint]*models.Wallet
	lastID  uint
}

func (m *mockWalletRepo) GetForUser(userID uint) ([]models.Wallet, error) {
	var list []models.Wallet
	for _, w := range m.wallets {
		if w.UserID == userID {
			list = append(list, *w)
		}
	}
	return list, nil
}

func (m *mockWalletRepo) CountForUser(userID uint) (int64, error) {
	var count int64
	for _, w := range m.wallets {
		if w.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *mockWalletRepo) GetByID(id uint) (*models.Wallet, error) {
	w, ok := m.wallets[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return w, nil
}

func (m *mockWalletRepo) GetByIDWithLock(tx *gorm.DB, id uint) (*models.Wallet, error) {
	return m.GetByID(id)
}

func (m *mockWalletRepo) Create(w *models.Wallet) error {
	m.lastID++
	w.ID = m.lastID
	m.wallets[w.ID] = w
	return nil
}

func (m *mockWalletRepo) Update(w *models.Wallet) error {
	m.wallets[w.ID] = w
	return nil
}

func (m *mockWalletRepo) UpdateWithTx(tx *gorm.DB, w *models.Wallet) error {
	return m.Update(w)
}

func (m *mockWalletRepo) Delete(id uint) error {
	delete(m.wallets, id)
	return nil
}

func TestWalletService_Create_LimitFreeUsers(t *testing.T) {
	userRepo := &mockUserRepo{users: make(map[string]*models.User)}
	walletRepo := &mockWalletRepo{wallets: make(map[uint]*models.Wallet)}
	srv := NewWalletService(walletRepo, userRepo)

	// Setup Free User
	freeUser := &models.User{
		Name:  "Free User",
		Email: "free@example.com",
		IsPro: false,
	}
	userRepo.Create(freeUser)

	// Create 1st Wallet (Success)
	_, err := srv.Create(freeUser.ID, "Cash", "cash", 10000, 1, 2)
	if err != nil {
		t.Fatalf("failed to create 1st wallet: %v", err)
	}

	// Create 2nd Wallet (Success)
	_, err = srv.Create(freeUser.ID, "Bank", "bank", 50000, 2, 3)
	if err != nil {
		t.Fatalf("failed to create 2nd wallet: %v", err)
	}

	// Create 3rd Wallet (Should Fail)
	_, err = srv.Create(freeUser.ID, "E-Wallet", "ewallet", 10000, 3, 4)
	if err == nil {
		t.Errorf("expected creating 3rd wallet for Free User to fail due to wallet limit")
	}

	// Setup Pro User
	proUser := &models.User{
		Name:  "Pro User",
		Email: "pro@example.com",
		IsPro: true,
	}
	userRepo.Create(proUser)

	// Create 3 wallets for Pro User (Should Success)
	for i := 1; i <= 3; i++ {
		_, err := srv.Create(proUser.ID, "Wallet", "cash", 1000, 1, 2)
		if err != nil {
			t.Fatalf("failed to create wallet %d for Pro User: %v", i, err)
		}
	}
}
