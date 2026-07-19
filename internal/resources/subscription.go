package resources

import (
	"finance-app-backend/internal/models"
)

type SubscriptionLogResource struct {
	ID             uint    `json:"id"`
	Plan           string  `json:"plan"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	PaymentMethod  *string `json:"payment_method"`
	PaymentChannel *string `json:"payment_channel"`
	CheckoutUrl    *string `json:"checkout_url"`
	XenditID       *string `json:"xendit_id"`
	StartsAt       string  `json:"starts_at"`
	EndsAt         string  `json:"ends_at"`
	CreatedAt      string  `json:"created_at"`
}

func ToSubscriptionLogResource(s *models.SubscriptionLog) SubscriptionLogResource {
	plan := ""
	if s.PlanID != nil {
		plan = *s.PlanID
	}

	return SubscriptionLogResource{
		ID:             s.ID,
		Plan:           plan,
		Amount:         s.Amount,
		Status:         s.Status,
		PaymentMethod:  s.PaymentMethod,
		PaymentChannel: s.PaymentChannel,
		CheckoutUrl:    s.XenditInvoiceURL,
		XenditID:       s.XenditInvoiceID,
		StartsAt:       s.StartsAt.Format("2006-01-02 15:04:05"),
		EndsAt:         s.EndsAt.Format("2006-01-02 15:04:05"),
		CreatedAt:      s.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToSubscriptionLogCollection(logs []models.SubscriptionLog) []SubscriptionLogResource {
	res := make([]SubscriptionLogResource, len(logs))
	for i, l := range logs {
		res[i] = ToSubscriptionLogResource(&l)
	}
	return res
}
