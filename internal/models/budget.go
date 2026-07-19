package models

import (
	"time"
)

// Budget represents the budgets table in database
type Budget struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;uniqueIndex:idx_user_category_month_year" json:"user_id"`
	CategoryID uint      `gorm:"not null;uniqueIndex:idx_user_category_month_year" json:"category_id"`
	Amount     float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Month      int       `gorm:"type:smallint;not null;uniqueIndex:idx_user_category_month_year" json:"month"`
	Year       int       `gorm:"type:smallint;not null;uniqueIndex:idx_user_category_month_year" json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Virtual field for response
	Spent float64 `gorm:"-" json:"spent"`

	// Relationships
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
