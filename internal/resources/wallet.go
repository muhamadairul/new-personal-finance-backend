package resources

import (
	"finance-app-backend/internal/models"
)

// WalletResource maps Wallet model to API output
type WalletResource struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
	Icon    int     `json:"icon"`
	Color   int64   `json:"color"`
}

// ToWalletResource transforms Wallet to resource
func ToWalletResource(w *models.Wallet) WalletResource {
	return WalletResource{
		ID:      w.ID,
		Name:    w.Name,
		Type:    w.Type,
		Balance: w.Balance,
		Icon:    w.Icon,
		Color:   w.Color,
	}
}

// ToWalletCollection transforms a slice of wallets
func ToWalletCollection(wallets []models.Wallet) []WalletResource {
	res := make([]WalletResource, len(wallets))
	for i, w := range wallets {
		res[i] = ToWalletResource(&w)
	}
	return res
}
