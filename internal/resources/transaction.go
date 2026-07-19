package resources

import (
	"finance-app-backend/internal/models"
)

// TransactionResource maps Transaction model to API response
type TransactionResource struct {
	ID            uint    `json:"id"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	CategoryID    uint    `json:"category_id"`
	WalletID      uint    `json:"wallet_id"`
	Note          *string `json:"note"`
	Date          string  `json:"date"`
	CreatedAt     string  `json:"created_at"`
	CategoryName  *string `json:"category_name"`
	CategoryIcon  *int    `json:"category_icon"`
	CategoryColor *int64  `json:"category_color"`
	WalletName    *string `json:"wallet_name"`
}

// ToTransactionResource transforms Transaction to resource
func ToTransactionResource(t *models.Transaction) TransactionResource {
	var catName *string
	var catIcon *int
	var catColor *int64
	if t.Category != nil {
		catName = &t.Category.Name
		catIcon = &t.Category.Icon
		catColor = &t.Category.Color
	}

	var wName *string
	if t.Wallet != nil {
		wName = &t.Wallet.Name
	}

	return TransactionResource{
		ID:            t.ID,
		Type:          t.Type,
		Amount:        t.Amount,
		CategoryID:    t.CategoryID,
		WalletID:      t.WalletID,
		Note:          t.Note,
		Date:          t.Date.Format("2006-01-02 15:04:05"),
		CreatedAt:     t.CreatedAt.Format("2006-01-02T15:04:05.000000Z"),
		CategoryName:  catName,
		CategoryIcon:  catIcon,
		CategoryColor: catColor,
		WalletName:    wName,
	}
}

// ToTransactionCollection transforms a slice of transactions
func ToTransactionCollection(txs []models.Transaction) []TransactionResource {
	res := make([]TransactionResource, len(txs))
	for i, t := range txs {
		res[i] = ToTransactionResource(&t)
	}
	return res
}
