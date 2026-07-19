package services

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"finance-app-backend/internal/config"
	"finance-app-backend/internal/models"
	"finance-app-backend/internal/pkg/mailer"
	"finance-app-backend/internal/pkg/utils"
	"finance-app-backend/internal/repositories"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  repositories.UserRepositoryInterface
	resetRepo repositories.PasswordResetRepositoryInterface
}

func NewAuthService(userRepo repositories.UserRepositoryInterface, resetRepo repositories.PasswordResetRepositoryInterface) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		resetRepo: resetRepo,
	}
}

func (s *AuthService) Register(name, email, password string) (*models.User, string, error) {
	_, err := s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, "", errors.New("email sudah terdaftar")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: &hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID, config.AppConfig.JwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("email atau kata sandi salah")
	}

	if user.Password == nil || !utils.CheckPasswordHash(password, *user.Password) {
		return nil, "", errors.New("email atau kata sandi salah")
	}

	token, err := utils.GenerateToken(user.ID, config.AppConfig.JwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) SocialLogin(provider, idToken, name, email string) (*models.User, string, error) {
	if provider != "google" {
		return nil, "", errors.New("provider tidak didukung")
	}

	// Verify token with Google API
	googleUserID, err := s.verifyGoogleIdToken(idToken)
	if err != nil {
		return nil, "", fmt.Errorf("autentikasi Google gagal: %w", err)
	}

	user, err := s.userRepo.GetByProvider(provider, googleUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Find by email next to check if user already registered via password
			existingUser, emailErr := s.userRepo.GetByEmail(email)
			if emailErr == nil {
				// Link provider
				existingUser.Provider = &provider
				existingUser.ProviderID = &googleUserID
				if err := s.userRepo.Update(existingUser); err != nil {
					return nil, "", err
				}
				user = existingUser
			} else {
				// Create new user without password
				user = &models.User{
					Name:       name,
					Email:      email,
					Provider:   &provider,
					ProviderID: &googleUserID,
				}
				if err := s.userRepo.Create(user); err != nil {
					return nil, "", err
				}
			}
		} else {
			return nil, "", err
		}
	}

	token, err := utils.GenerateToken(user.ID, config.AppConfig.JwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) verifyGoogleIdToken(idToken string) (string, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token verification status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	// Check audience matches expected client IDs
	aud, _ := data["aud"].(string)
	clientGoogleID := os.Getenv("GOOGLE_CLIENT_ID")
	clientGoogleAndroidID := os.Getenv("GOOGLE_CLIENT_ID_ANDROID")

	if aud != clientGoogleID && aud != clientGoogleAndroidID {
		return "", errors.New("google token audience mismatch")
	}

	sub, _ := data["sub"].(string)
	if sub == "" {
		return "", errors.New("google token missing sub claim")
	}

	return sub, nil
}

func (s *AuthService) SendOtp(email string) error {
	// Verify user exists first
	_, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return errors.New("kami tidak dapat menemukan pengguna dengan alamat email tersebut")
	}

	// Generate 6 digit OTP
	otpNum, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return err
	}
	otp := fmt.Sprintf("%06d", otpNum.Int64()+100000)

	// Clean up old resets
	s.resetRepo.DeleteByEmail(email)

	// Hash OTP for secure storage
	hashedOtp, err := utils.HashPassword(otp)
	if err != nil {
		return err
	}

	now := time.Now()
	resetToken := &models.PasswordResetToken{
		Email:     email,
		Token:     hashedOtp,
		CreatedAt: &now,
	}

	if err := s.resetRepo.Create(resetToken); err != nil {
		return err
	}

	// Send HTML email
	return mailer.SendOtpEmail(email, otp)
}

func (s *AuthService) VerifyOtp(email, otp string) error {
	record, err := s.resetRepo.GetByEmail(email)
	if err != nil {
		return errors.New("kode OTP tidak valid atau salah")
	}

	if !utils.CheckPasswordHash(otp, record.Token) {
		return errors.New("kode OTP tidak valid atau salah")
	}

	// Check expiration (15 minutes)
	if record.CreatedAt != nil && record.CreatedAt.Add(15*time.Minute).Before(time.Now()) {
		return errors.New("kode OTP sudah kadaluarsa")
	}

	return nil
}

func (s *AuthService) ResetPassword(email, otp, newPassword string) error {
	if err := s.VerifyOtp(email, otp); err != nil {
		return err
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = &hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Invalidate OTP
	s.resetRepo.DeleteByEmail(email)

	return nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *AuthService) UpdateProfile(userID uint, name string, phone, address, dateOfBirth, gender *string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.Phone = phone
	user.Address = address

	if dateOfBirth != nil && *dateOfBirth != "" {
		t, err := time.Parse("2006-01-02", *dateOfBirth)
		if err != nil {
			return nil, errors.New("format tanggal lahir harus YYYY-MM-DD")
		}
		user.DateOfBirth = &t
	} else {
		user.DateOfBirth = nil
	}

	user.Gender = gender

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) UpdateFcmToken(userID uint, fcmToken string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.FcmToken = &fcmToken
	return s.userRepo.Update(user)
}

func (s *AuthService) ClearFcmTokenAndLogout(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.FcmToken = nil
	return s.userRepo.Update(user)
}

func (s *AuthService) SavePhoto(userID uint, fileHeader io.Reader, filename string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Delete old photo if exists locally
	if user.PhotoURL != nil {
		oldPath := filepath.Join("public", "storage", filepath.Base(*user.PhotoURL))
		os.Remove(oldPath)
	}

	// Make directory public/storage if not exists
	storageDir := filepath.Join("public", "storage")
	os.MkdirAll(storageDir, os.ModePerm)

	// Save new photo
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("profile_%d_%d%s", userID, time.Now().Unix(), ext)
	dstPath := filepath.Join(storageDir, newFilename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, fileHeader); err != nil {
		return nil, err
	}

	photoPath := filepath.Join("profile-photos", newFilename)
	user.PhotoURL = &photoPath

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) DeletePhoto(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user.PhotoURL != nil {
		oldPath := filepath.Join("public", "storage", filepath.Base(*user.PhotoURL))
		os.Remove(oldPath)
	}

	user.PhotoURL = nil
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
