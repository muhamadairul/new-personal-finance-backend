package models

import (
	"time"
)

// User represents the users table in database
type User struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	Name              string            `gorm:"type:varchar(255);not null" json:"name"`
	Email             string            `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	EmailVerifiedAt   *time.Time        `json:"email_verified_at,omitempty"`
	Password          *string           `gorm:"type:varchar(255)" json:"-"`
	Provider          *string           `gorm:"type:varchar(20);index:idx_provider_id" json:"provider,omitempty"`
	ProviderID        *string           `gorm:"type:varchar(255);index:idx_provider_id" json:"provider_id,omitempty"`
	PhotoURL          *string           `gorm:"type:varchar(500)" json:"photo_url"`
	Phone             *string           `gorm:"type:varchar(20)" json:"phone"`
	Address           *string           `gorm:"type:text" json:"address"`
	DateOfBirth       *time.Time        `gorm:"type:date" json:"date_of_birth"`
	Gender            *string           `gorm:"type:varchar(1)" json:"gender"`
	IsPro             bool              `gorm:"default:false" json:"is_pro"`
	SubscriptionUntil *time.Time        `json:"subscription_until"`
	FcmToken          *string           `gorm:"type:varchar(255)" json:"fcm_token,omitempty"`
	IsAdmin           bool              `gorm:"default:false" json:"is_admin"`
	RememberToken     *string           `gorm:"type:varchar(100)" json:"-"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`

	// Relationships
	Wallets          []Wallet          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"wallets,omitempty"`
	Transactions     []Transaction     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"transactions,omitempty"`
	Categories       []Category        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"categories,omitempty"`
	Budgets          []Budget          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"budgets,omitempty"`
	SubscriptionLogs []SubscriptionLog `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"subscription_logs,omitempty"`
}

// CheckIsPro returns whether user currently has active pro subscription
func (u *User) CheckIsPro() bool {
	if !u.IsPro {
		return false
	}
	if u.SubscriptionUntil != nil {
		return u.SubscriptionUntil.After(time.Now())
	}
	return true
}
