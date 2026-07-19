package repositories

import (
	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

// CategoryRepositoryInterface specifies category DB operations
type CategoryRepositoryInterface interface {
	GetForUser(userID uint) ([]models.Category, error)
	GetByID(id uint) (*models.Category, error)
	Create(cat *models.Category) error
	Update(cat *models.Category) error
	Delete(id uint) error
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetForUser(userID uint) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("user_id IS NULL OR user_id = ?", userID).
		Order("type ASC").
		Order("id ASC").
		Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(cat *models.Category) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepository) Update(cat *models.Category) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
