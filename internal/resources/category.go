package resources

import (
	"finance-app-backend/internal/models"
)

// CategoryResource maps the Category model to API response
type CategoryResource struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Icon  int    `json:"icon"`
	Color int64  `json:"color"`
	Type  string `json:"type"`
}

// ToCategoryResource transforms a GORM Category model to CategoryResource
func ToCategoryResource(cat *models.Category) CategoryResource {
	return CategoryResource{
		ID:    cat.ID,
		Name:  cat.Name,
		Icon:  cat.Icon,
		Color: cat.Color,
		Type:  cat.Type,
	}
}

// ToCategoryCollection transforms a slice of Category models into resources
func ToCategoryCollection(cats []models.Category) []CategoryResource {
	res := make([]CategoryResource, len(cats))
	for i, c := range cats {
		res[i] = ToCategoryResource(&c)
	}
	return res
}
