package services

import (
	"errors"
	"fmt"
	"time"

	"finance-app-backend/internal/models"
	"finance-app-backend/internal/repositories"

	"gorm.io/gorm"
)

type TransactionService struct {
	txRepo       repositories.TransactionRepositoryInterface
	walletRepo   repositories.WalletRepositoryInterface
	categoryRepo repositories.CategoryRepositoryInterface
	db           *gorm.DB
}

func NewTransactionService(
	txRepo repositories.TransactionRepositoryInterface,
	walletRepo repositories.WalletRepositoryInterface,
	categoryRepo repositories.CategoryRepositoryInterface,
	db *gorm.DB,
) *TransactionService {
	return &TransactionService{
		txRepo:       txRepo,
		walletRepo:   walletRepo,
		categoryRepo: categoryRepo,
		db:           db,
	}
}

func (s *TransactionService) List(userID uint, month, year int) ([]models.Transaction, error) {
	return s.txRepo.GetForUser(userID, month, year)
}

func (s *TransactionService) GetByID(userID uint, id uint) (*models.Transaction, error) {
	t, err := s.txRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if t.UserID != userID {
		return nil, errors.New("akses ditolak")
	}

	return t, nil
}

func (s *TransactionService) Create(
	userID uint,
	typeStr string,
	amount float64,
	categoryID, walletID uint,
	note *string,
	date time.Time,
) (*models.Transaction, error) {
	// Verify category belongs to user (or is global)
	cat, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if cat.UserID != nil && *cat.UserID != userID {
		return nil, errors.New("akses kategori ditolak")
	}

	var transaction *models.Transaction

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Lock wallet for update
		wallet, err := s.walletRepo.GetByIDWithLock(tx, walletID)
		if err != nil {
			return errors.New("dompet tidak ditemukan")
		}
		if wallet.UserID != userID {
			return errors.New("akses dompet ditolak")
		}

		// Validate sufficient balance for expense
		if typeStr == "expense" && wallet.Balance < amount {
			return fmt.Errorf("Saldo dompet tidak mencukupi untuk transaksi ini. Saldo tersedia: Rp %.0f", wallet.Balance)
		}

		transaction = &models.Transaction{
			UserID:     userID,
			WalletID:   walletID,
			CategoryID: categoryID,
			Type:       typeStr,
			Amount:     amount,
			Note:       note,
			Date:       date,
		}

		if err := s.txRepo.Create(tx, transaction); err != nil {
			return err
		}

		// Update wallet balance
		if typeStr == "income" {
			wallet.Balance += amount
		} else {
			wallet.Balance -= amount
		}

		if err := s.walletRepo.UpdateWithTx(tx, wallet); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload relations
	return s.txRepo.GetByID(transaction.ID)
}

func (s *TransactionService) Update(
	userID uint,
	id uint,
	typeStr string,
	amount float64,
	categoryID, walletID uint,
	note *string,
	date time.Time,
) (*models.Transaction, error) {
	// Verify category
	cat, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if cat.UserID != nil && *cat.UserID != userID {
		return nil, errors.New("akses kategori ditolak")
	}

	var transaction *models.Transaction

	err = s.db.Transaction(func(tx *gorm.DB) error {
		t, err := s.txRepo.GetByID(id)
		if err != nil {
			return errors.New("transaksi tidak ditemukan")
		}
		if t.UserID != userID {
			return errors.New("akses transaksi ditolak")
		}

		originalAmount := t.Amount
		originalType := t.Type
		originalWalletID := t.WalletID

		// Revert old wallet balance
		oldWallet, err := s.walletRepo.GetByIDWithLock(tx, originalWalletID)
		if err != nil {
			return errors.New("dompet asli tidak ditemukan")
		}
		if originalType == "income" {
			oldWallet.Balance -= originalAmount
		} else {
			oldWallet.Balance += originalAmount
		}
		if err := s.walletRepo.UpdateWithTx(tx, oldWallet); err != nil {
			return err
		}

		// Fetch and lock new wallet
		var newWallet *models.Wallet
		if originalWalletID == walletID {
			newWallet = oldWallet
		} else {
			newWallet, err = s.walletRepo.GetByIDWithLock(tx, walletID)
			if err != nil {
				return errors.New("dompet baru tidak ditemukan")
			}
			if newWallet.UserID != userID {
				return errors.New("akses dompet baru ditolak")
			}
		}

		// Validate sufficient balance in new wallet
		if typeStr == "expense" && newWallet.Balance < amount {
			return fmt.Errorf("Saldo dompet tidak mencukupi untuk transaksi ini. Saldo tersedia: Rp %.0f", newWallet.Balance)
		}

		// Apply new balance
		if typeStr == "income" {
			newWallet.Balance += amount
		} else {
			newWallet.Balance -= amount
		}
		if err := s.walletRepo.UpdateWithTx(tx, newWallet); err != nil {
			return err
		}

		t.Type = typeStr
		t.Amount = amount
		t.CategoryID = categoryID
		t.WalletID = walletID
		t.Note = note
		t.Date = date

		if err := s.txRepo.Update(tx, t); err != nil {
			return err
		}

		transaction = t
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.txRepo.GetByID(transaction.ID)
}

func (s *TransactionService) Delete(userID uint, id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		t, err := s.txRepo.GetByID(id)
		if err != nil {
			return errors.New("transaksi tidak ditemukan")
		}
		if t.UserID != userID {
			return errors.New("akses transaksi ditolak")
		}

		// Revert wallet balance
		wallet, err := s.walletRepo.GetByIDWithLock(tx, t.WalletID)
		if err != nil {
			return errors.New("dompet tidak ditemukan")
		}
		if t.Type == "income" {
			wallet.Balance -= t.Amount
		} else {
			wallet.Balance += t.Amount
		}
		if err := s.walletRepo.UpdateWithTx(tx, wallet); err != nil {
			return err
		}

		return s.txRepo.Delete(tx, id)
	})
}
