package models

import (
	"time"
)

// Transaction represents the transactions table in database
type Transaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index:idx_user_date;index:idx_user_type" json:"user_id"`
	WalletID   uint      `gorm:"not null;index" json:"wallet_id"`
	CategoryID uint      `gorm:"not null;index" json:"category_id"`
	Type       string    `gorm:"type:varchar(50);not null" json:"type"`
	Amount     float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Note       *string   `gorm:"type:text" json:"note"`
	Date       time.Time `gorm:"not null;index:idx_user_date" json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relationships
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Wallet   *Wallet   `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
