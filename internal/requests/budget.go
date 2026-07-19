package requests

// BudgetRequest holds budget payload validation
type BudgetRequest struct {
	CategoryID uint    `json:"category_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,numeric,gt=0"`
	Month      int     `json:"month" validate:"required,min=1,max=12"`
	Year       int     `json:"year" validate:"required,min=2000,max=2100"`
}
