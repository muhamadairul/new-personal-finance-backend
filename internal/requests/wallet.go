package requests

// WalletRequest holds wallet payload validation
type WalletRequest struct {
	Name    string   `json:"name" validate:"required,max=255"`
	Type    string   `json:"type" validate:"required,oneof=cash bank ewallet"`
	Balance *float64 `json:"balance" validate:"omitempty,numeric,min=0"`
	Icon    int      `json:"icon" validate:"required"`
	Color   int64    `json:"color" validate:"required"`
}
