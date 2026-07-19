package requests

// TransactionRequest holds transaction payload validation
type TransactionRequest struct {
	Type       string  `json:"type" validate:"required,oneof=income expense"`
	Amount     float64 `json:"amount" validate:"required,numeric,gt=0"`
	CategoryID uint    `json:"category_id" validate:"required"`
	WalletID   uint    `json:"wallet_id" validate:"required"`
	Note       *string `json:"note" validate:"omitempty,max=500"`
	Date       string  `json:"date" validate:"required"` // Expected date string
}
