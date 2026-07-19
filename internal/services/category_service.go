package services

import (
	"errors"

	"finance-app-backend/internal/models"
	"finance-app-backend/internal/repositories"
)

type CategoryService struct {
	categoryRepo repositories.CategoryRepositoryInterface
	userRepo     repositories.UserRepositoryInterface
}

func NewCategoryService(categoryRepo repositories.CategoryRepositoryInterface, userRepo repositories.UserRepositoryInterface) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}
}

func (s *CategoryService) List(userID uint) ([]models.Category, error) {
	return s.categoryRepo.GetForUser(userID)
}

func (s *CategoryService) Create(userID uint, name string, icon int, color int64, cType string) (*models.Category, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.CheckIsPro() {
		return nil, errors.New("pengguna gratis hanya bisa menggunakan kategori default")
	}

	cat := &models.Category{
		UserID: &userID,
		Name:   name,
		Icon:   icon,
		Color:  color,
		Type:   cType,
	}

	if err := s.categoryRepo.Create(cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *CategoryService) GetByID(userID uint, id uint) (*models.Category, error) {
	cat, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify ownership: global or owned by requester
	if cat.UserID != nil && *cat.UserID != userID {
		return nil, errors.New("akses ditolak")
	}

	return cat, nil
}

func (s *CategoryService) Update(userID uint, id uint, name string, icon int, color int64, cType string) (*models.Category, error) {
	cat, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if cat.UserID == nil || *cat.UserID != userID {
		return nil, errors.New("tidak dapat mengubah kategori default")
	}

	cat.Name = name
	cat.Icon = icon
	cat.Color = color
	cat.Type = cType

	if err := s.categoryRepo.Update(cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *CategoryService) Delete(userID uint, id uint) error {
	cat, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return err
	}

	if cat.UserID == nil || *cat.UserID != userID {
		return errors.New("tidak dapat menghapus kategori default")
	}

	return s.categoryRepo.Delete(id)
}
