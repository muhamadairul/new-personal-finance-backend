package services

import (
	"errors"

	"finance-app-backend/internal/models"
	"finance-app-backend/internal/repositories"

	"gorm.io/gorm"
)

type BudgetService struct {
	budgetRepo   repositories.BudgetRepositoryInterface
	categoryRepo repositories.CategoryRepositoryInterface
	db           *gorm.DB
}

func NewBudgetService(
	budgetRepo repositories.BudgetRepositoryInterface,
	categoryRepo repositories.CategoryRepositoryInterface,
	db *gorm.DB,
) *BudgetService {
	return &BudgetService{
		budgetRepo:   budgetRepo,
		categoryRepo: categoryRepo,
		db:           db,
	}
}

func (s *BudgetService) List(userID uint, month, year int) ([]models.Budget, error) {
	budgets, err := s.budgetRepo.GetForUser(userID, month, year)
	if err != nil {
		return nil, err
	}

	// Calculate spent for each budget
	for i := range budgets {
		spent, err := s.budgetRepo.GetSpentAmount(userID, budgets[i].CategoryID, budgets[i].Month, budgets[i].Year)
		if err == nil {
			budgets[i].Spent = spent
		}
	}

	return budgets, nil
}

func (s *BudgetService) Upsert(userID, categoryID uint, amount float64, month, year int) (*models.Budget, error) {
	// Verify category
	cat, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if cat.UserID != nil && *cat.UserID != userID {
		return nil, errors.New("akses kategori ditolak")
	}

	var budget models.Budget
	err = s.db.Where("user_id = ? AND category_id = ? AND month = ? AND year = ?", userID, categoryID, month, year).
		First(&budget).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		budget = models.Budget{
			UserID:     userID,
			CategoryID: categoryID,
			Amount:     amount,
			Month:      month,
			Year:       year,
		}
		if err := s.budgetRepo.Upsert(&budget); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		budget.Amount = amount
		if err := s.budgetRepo.Update(&budget); err != nil {
			return nil, err
		}
	}

	// Reload relations and calculate spent
	res, err := s.budgetRepo.GetByID(budget.ID)
	if err == nil {
		spent, _ := s.budgetRepo.GetSpentAmount(userID, categoryID, month, year)
		res.Spent = spent
	}
	return res, err
}

func (s *BudgetService) GetByID(userID uint, id uint) (*models.Budget, error) {
	b, err := s.budgetRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if b.UserID != userID {
		return nil, errors.New("akses ditolak")
	}

	spent, _ := s.budgetRepo.GetSpentAmount(userID, b.CategoryID, b.Month, b.Year)
	b.Spent = spent

	return b, nil
}

func (s *BudgetService) Update(userID uint, id uint, categoryID uint, amount float64, month, year int) (*models.Budget, error) {
	b, err := s.budgetRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if b.UserID != userID {
		return nil, errors.New("akses ditolak")
	}

	// Verify new category
	cat, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if cat.UserID != nil && *cat.UserID != userID {
		return nil, errors.New("akses kategori ditolak")
	}

	b.CategoryID = categoryID
	b.Amount = amount
	b.Month = month
	b.Year = year

	if err := s.budgetRepo.Update(b); err != nil {
		return nil, err
	}

	spent, _ := s.budgetRepo.GetSpentAmount(userID, categoryID, month, year)
	b.Spent = spent

	return b, nil
}

func (s *BudgetService) Delete(userID uint, id uint) error {
	b, err := s.budgetRepo.GetByID(id)
	if err != nil {
		return err
	}

	if b.UserID != userID {
		return errors.New("akses ditolak")
	}

	return s.budgetRepo.Delete(id)
}
