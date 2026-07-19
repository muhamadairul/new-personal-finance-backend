package services

import (
	"errors"

	"finance-app-backend/internal/models"
	"finance-app-backend/internal/repositories"
)

type WalletService struct {
	walletRepo repositories.WalletRepositoryInterface
	userRepo   repositories.UserRepositoryInterface
}

func NewWalletService(walletRepo repositories.WalletRepositoryInterface, userRepo repositories.UserRepositoryInterface) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		userRepo:   userRepo,
	}
}

func (s *WalletService) List(userID uint) ([]models.Wallet, error) {
	return s.walletRepo.GetForUser(userID)
}

func (s *WalletService) Create(userID uint, name, wType string, balance float64, icon int, color int64) (*models.Wallet, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.CheckIsPro() {
		count, err := s.walletRepo.CountForUser(userID)
		if err != nil {
			return nil, err
		}
		if count >= 2 {
			return nil, errors.New("pengguna gratis hanya bisa memiliki maksimal 2 dompet")
		}
	}

	wallet := &models.Wallet{
		UserID:  userID,
		Name:    name,
		Type:    wType,
		Balance: balance,
		Icon:    icon,
		Color:   color,
	}

	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) GetByID(userID uint, id uint) (*models.Wallet, error) {
	wallet, err := s.walletRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if wallet.UserID != userID {
		return nil, errors.New("akses ditolak")
	}

	return wallet, nil
}

func (s *WalletService) Update(userID uint, id uint, name, wType string, balance float64, icon int, color int64) (*models.Wallet, error) {
	wallet, err := s.walletRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if wallet.UserID != userID {
		return nil, errors.New("akses ditolak")
	}

	wallet.Name = name
	wallet.Type = wType
	wallet.Balance = balance
	wallet.Icon = icon
	wallet.Color = color

	if err := s.walletRepo.Update(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) Delete(userID uint, id uint) error {
	wallet, err := s.walletRepo.GetByID(id)
	if err != nil {
		return err
	}

	if wallet.UserID != userID {
		return errors.New("akses ditolak")
	}

	return s.walletRepo.Delete(id)
}
