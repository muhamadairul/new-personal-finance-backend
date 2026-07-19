package seeders

import (
	"log"

	"finance-app-backend/internal/models"

	"gorm.io/gorm"
)

type CategorySeeder struct{}

func (s *CategorySeeder) Name() string {
	return "CategorySeeder"
}

func (s *CategorySeeder) Run(db *gorm.DB) error {
	categories := []models.Category{
		// Expense categories (matching Flutter DefaultCategories)
		{ID: 1, Name: "Makan", Icon: 0xe532, Color: 0xFFFF6B6B, Type: "expense", UserID: nil},
		{ID: 2, Name: "Transportasi", Icon: 0xe1d7, Color: 0xFF4ECDC4, Type: "expense", UserID: nil},
		{ID: 3, Name: "Belanja", Icon: 0xf37c, Color: 0xFFFFBE0B, Type: "expense", UserID: nil},
		{ID: 4, Name: "Tagihan", Icon: 0xe4c0, Color: 0xFF845EC2, Type: "expense", UserID: nil},
		{ID: 5, Name: "Hiburan", Icon: 0xe40c, Color: 0xFFFF9671, Type: "expense", UserID: nil},
		{ID: 6, Name: "Kesehatan", Icon: 0xf109, Color: 0xFF00C9A7, Type: "expense", UserID: nil},
		{ID: 7, Name: "Pendidikan", Icon: 0xe559, Color: 0xFF4D8076, Type: "expense", UserID: nil},
		{ID: 8, Name: "Lainnya", Icon: 0xe400, Color: 0xFF8E8E93, Type: "expense", UserID: nil},

		// Income categories
		{ID: 9, Name: "Gaji", Icon: 0xe850, Color: 0xFF00C853, Type: "income", UserID: nil},
		{ID: 10, Name: "Freelance", Icon: 0xe3e9, Color: 0xFF2196F3, Type: "income", UserID: nil},
		{ID: 11, Name: "Investasi", Icon: 0xe8e5, Color: 0xFFFF9800, Type: "income", UserID: nil},
		{ID: 12, Name: "Hadiah", Icon: 0xe8f6, Color: 0xFFE91E63, Type: "income", UserID: nil},
		{ID: 13, Name: "Lainnya", Icon: 0xe400, Color: 0xFF8E8E93, Type: "income", UserID: nil},
	}

	for _, cat := range categories {
		var existing models.Category
		err := db.Where("id = ?", cat.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&cat).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			log.Printf("[CategorySeeder] Category ID %d (%s) already exists, skipping.", cat.ID, cat.Name)
		}
	}

	return nil
}

func init() {
	Register(&CategorySeeder{})
}
