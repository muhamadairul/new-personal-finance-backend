package repositories

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

type BudgetRepositoryInterface interface {
	GetForUser(userID uint, month, year int) ([]models.Budget, error)
	GetByID(id uint) (*models.Budget, error)
	GetSpentAmount(userID, categoryID uint, month, year int) (float64, error)
	Upsert(b *models.Budget) error
	Update(b *models.Budget) error
	Delete(id uint) error
}

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) GetForUser(userID uint, month, year int) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.db.Where("user_id = ? AND month = ? AND year = ?", userID, month, year).
		Preload("Category").
		Find(&budgets).Error
	return budgets, err
}

func (r *BudgetRepository) GetByID(id uint) (*models.Budget, error) {
	var b models.Budget
	err := r.db.Preload("Category").First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BudgetRepository) GetSpentAmount(userID, categoryID uint, month, year int) (float64, error) {
	var sum float64
	row := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND category_id = ? AND type = 'expense' AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", userID, categoryID, month, year).
		Row()
	err := row.Scan(&sum)
	return sum, err
}

func (r *BudgetRepository) Upsert(b *models.Budget) error {
	// GORM Clause OnConflict for upserting
	return r.db.Save(b).Error
}

func (r *BudgetRepository) Update(b *models.Budget) error {
	return r.db.Save(b).Error
}

func (r *BudgetRepository) Delete(id uint) error {
	return r.db.Delete(&models.Budget{}, id).Error
}
