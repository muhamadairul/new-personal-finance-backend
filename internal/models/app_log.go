package models

import (
	"time"
)

// AppLog represents the app_logs table in database
type AppLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Level     string    `gorm:"type:varchar(20);not null;index" json:"level"`
	Channel   string    `gorm:"type:varchar(50);default:'stack';not null" json:"channel"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Context   *string   `gorm:"type:text" json:"context,omitempty"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
