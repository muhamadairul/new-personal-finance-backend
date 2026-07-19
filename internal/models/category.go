package models

import (
	"time"
)

// Category represents the categories table in database
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Icon      int       `gorm:"not null" json:"icon"`
	Color     int64     `gorm:"not null" json:"color"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"transactions,omitempty"`
	Budgets      []Budget      `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"budgets,omitempty"`
}

// IsGlobal checks if category is global default category
func (c *Category) IsGlobal() bool {
	return c.UserID == nil
}
