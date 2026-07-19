package services

import (
	"errors"
	"testing"

	"finance-app-backend/internal/config"
	"finance-app-backend/internal/models"
	"finance-app-backend/internal/requests"
)

type mockUserRepo struct {
	users map[string]*models.User
	lastID uint
}

func (m *mockUserRepo) Create(user *models.User) error {
	m.lastID++
	user.ID = m.lastID
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByID(id uint) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) GetByEmail(email string) (*models.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetByProvider(provider, providerID string) (*models.User, error) {
	for _, u := range m.users {
		if u.Provider != nil && *u.Provider == provider && u.ProviderID != nil && *u.ProviderID == providerID {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) Update(user *models.User) error {
	m.users[user.Email] = user
	return nil
}

type mockResetRepo struct {
	tokens map[string]*models.PasswordResetToken
}

func (m *mockResetRepo) Create(token *models.PasswordResetToken) error {
	m.tokens[token.Email] = token
	return nil
}

func (m *mockResetRepo) GetByEmail(email string) (*models.PasswordResetToken, error) {
	t, ok := m.tokens[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}

func (m *mockResetRepo) DeleteByEmail(email string) error {
	delete(m.tokens, email)
	return nil
}

func TestAuthService_Register(t *testing.T) {
	config.AppConfig = &config.Config{
		JwtSecret: "test-secret",
	}

	userRepo := &mockUserRepo{users: make(map[string]*models.User)}
	resetRepo := &mockResetRepo{tokens: make(map[string]*models.PasswordResetToken)}
	srv := NewAuthService(userRepo, resetRepo)

	req := requests.RegisterRequest{
		Name:                 "Test User",
		Email:                "test@example.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	user, token, err := srv.Register(req.Name, req.Email, req.Password)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if user.Name != req.Name {
		t.Errorf("expected name %s, got %s", req.Name, user.Name)
	}

	if token == "" {
		t.Errorf("expected non-empty token")
	}

	// Try registering with same email
	_, _, err = srv.Register(req.Name, req.Email, req.Password)
	if err == nil {
		t.Fatalf("expected duplicate email registration to fail")
	}
}

func TestAuthService_Login(t *testing.T) {
	config.AppConfig = &config.Config{
		JwtSecret: "test-secret",
	}

	userRepo := &mockUserRepo{users: make(map[string]*models.User)}
	resetRepo := &mockResetRepo{tokens: make(map[string]*models.PasswordResetToken)}
	srv := NewAuthService(userRepo, resetRepo)

	// Seed user
	_, _, _ = srv.Register("Login User", "login@example.com", "password123")

	// Successful login
	user, token, err := srv.Login("login@example.com", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if user.Email != "login@example.com" {
		t.Errorf("expected email login@example.com, got %s", user.Email)
	}
	if token == "" {
		t.Errorf("expected non-empty token")
	}

	// Failed login (wrong password)
	_, _, err = srv.Login("login@example.com", "wrongpassword")
	if err == nil {
		t.Errorf("expected wrong password login to fail")
	}
}

func TestAuthService_UpdateProfile(t *testing.T) {
	config.AppConfig = &config.Config{
		JwtSecret: "test-secret",
	}

	userRepo := &mockUserRepo{users: make(map[string]*models.User)}
	resetRepo := &mockResetRepo{tokens: make(map[string]*models.PasswordResetToken)}
	srv := NewAuthService(userRepo, resetRepo)

	// Seed user
	user, _, _ := srv.Register("Old Name", "profile@example.com", "password123")

	phone := "08123456789"
	address := "New Address"
	dob := "1995-12-15"
	gender := "L"

	updatedUser, err := srv.UpdateProfile(user.ID, "New Name", &phone, &address, &dob, &gender)
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}

	if updatedUser.Name != "New Name" {
		t.Errorf("expected name New Name, got %s", updatedUser.Name)
	}
	if *updatedUser.Phone != phone {
		t.Errorf("expected phone %s, got %s", phone, *updatedUser.Phone)
	}
	if updatedUser.DateOfBirth == nil || updatedUser.DateOfBirth.Format("2006-01-02") != dob {
		t.Errorf("expected dob %s", dob)
	}
}
