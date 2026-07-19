package services

import (
	"errors"
	"testing"

	"finance-app-backend/internal/models"
)

type mockCategoryRepo struct {
	categories map[uint]*models.Category
	lastID     uint
}

func (m *mockCategoryRepo) GetForUser(userID uint) ([]models.Category, error) {
	var list []models.Category
	for _, c := range m.categories {
		if c.UserID == nil || *c.UserID == userID {
			list = append(list, *c)
		}
	}
	return list, nil
}

func (m *mockCategoryRepo) GetByID(id uint) (*models.Category, error) {
	c, ok := m.categories[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (m *mockCategoryRepo) Create(cat *models.Category) error {
	m.lastID++
	cat.ID = m.lastID
	m.categories[cat.ID] = cat
	return nil
}

func (m *mockCategoryRepo) Update(cat *models.Category) error {
	m.categories[cat.ID] = cat
	return nil
}

func (m *mockCategoryRepo) Delete(id uint) error {
	delete(m.categories, id)
	return nil
}

func TestCategoryService_Create_FreeVsPro(t *testing.T) {
	userRepo := &mockUserRepo{users: make(map[string]*models.User)}
	categoryRepo := &mockCategoryRepo{categories: make(map[uint]*models.Category)}
	srv := NewCategoryService(categoryRepo, userRepo)

	// 1. Setup Free User
	freeUser := &models.User{
		Name:  "Free User",
		Email: "free@example.com",
		IsPro: false,
	}
	userRepo.Create(freeUser)

	// Attempt creating category (should fail)
	_, err := srv.Create(freeUser.ID, "Kategori Baru", 1234, 5678, "expense")
	if err == nil {
		t.Errorf("expected Free User creating custom category to fail")
	}

	// 2. Setup Pro User
	proUser := &models.User{
		Name:  "Pro User",
		Email: "pro@example.com",
		IsPro: true,
	}
	userRepo.Create(proUser)

	// Attempt creating category (should succeed)
	cat, err := srv.Create(proUser.ID, "Kategori Baru", 1234, 5678, "expense")
	if err != nil {
		t.Fatalf("expected Pro User creating custom category to succeed: %v", err)
	}

	if cat.Name != "Kategori Baru" {
		t.Errorf("expected category name 'Kategori Baru', got %s", cat.Name)
	}
}

func TestCategoryService_Delete_GlobalSafety(t *testing.T) {
	userRepo := &mockUserRepo{users: make(map[string]*models.User)}
	categoryRepo := &mockCategoryRepo{categories: make(map[uint]*models.Category)}
	srv := NewCategoryService(categoryRepo, userRepo)

	// Create Global Category (UserID = nil)
	globalCat := &models.Category{
		Name:   "Global Cat",
		UserID: nil,
	}
	categoryRepo.Create(globalCat)

	// Try deleting global category (should fail)
	err := srv.Delete(1, globalCat.ID)
	if err == nil {
		t.Errorf("expected deleting global category to fail")
	}
}
