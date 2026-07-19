package models

import (
	"time"
)

// Wallet represents the wallets table in database
type Wallet struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Type      string    `gorm:"type:varchar(50);default:'cash';not null" json:"type"`
	Balance   float64   `gorm:"type:decimal(15,2);default:0;not null" json:"balance"`
	Icon      int       `gorm:"not null" json:"icon"`
	Color     int64     `gorm:"not null" json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE" json:"transactions,omitempty"`
}
