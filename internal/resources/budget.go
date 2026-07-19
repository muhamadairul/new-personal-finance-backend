package resources

import (
	"finance-app-backend/internal/models"
)

// BudgetResource maps Budget model to API response
type BudgetResource struct {
	ID            uint    `json:"id"`
	CategoryID    uint    `json:"category_id"`
	Amount        float64 `json:"amount"`
	Spent         float64 `json:"spent"`
	Month         int     `json:"month"`
	Year          int     `json:"year"`
	CategoryName  *string `json:"category_name"`
	CategoryIcon  *int    `json:"category_icon"`
	CategoryColor *int64  `json:"category_color"`
}

// ToBudgetResource transforms Budget to resource
func ToBudgetResource(b *models.Budget) BudgetResource {
	var catName *string
	var catIcon *int
	var catColor *int64
	if b.Category != nil {
		catName = &b.Category.Name
		catIcon = &b.Category.Icon
		catColor = &b.Category.Color
	}

	return BudgetResource{
		ID:            b.ID,
		CategoryID:    b.CategoryID,
		Amount:        b.Amount,
		Spent:         b.Spent,
		Month:         b.Month,
		Year:          b.Year,
		CategoryName:  catName,
		CategoryIcon:  catIcon,
		CategoryColor: catColor,
	}
}

// ToBudgetCollection transforms a slice of budgets
func ToBudgetCollection(budgets []models.Budget) []BudgetResource {
	res := make([]BudgetResource, len(budgets))
	for i, b := range budgets {
		res[i] = ToBudgetResource(&b)
	}
	return res
}
