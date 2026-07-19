package models

import (
	"time"
)

// SubscriptionLog represents the subscription_logs table in database
type SubscriptionLog struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	Type             string    `gorm:"type:varchar(50);not null" json:"type"`
	XenditInvoiceID  *string   `gorm:"type:varchar(255);index" json:"xendit_invoice_id,omitempty"`
	XenditInvoiceURL *string   `gorm:"type:varchar(500)" json:"xendit_invoice_url,omitempty"`
	PaymentMethod    *string   `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	PaymentChannel   *string   `gorm:"type:varchar(50)" json:"payment_channel,omitempty"`
	Status           string    `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	PlanID           *string   `gorm:"type:varchar(50)" json:"plan_id,omitempty"`
	Amount           float64   `gorm:"type:decimal(12,2);default:0;not null" json:"amount"`
	StartsAt         time.Time `gorm:"not null" json:"starts_at"`
	EndsAt           time.Time `gorm:"not null;index" json:"ends_at"`
	Notes            *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
