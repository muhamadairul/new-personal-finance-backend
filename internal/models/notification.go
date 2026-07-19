package models

import (
	"time"
)

// Notification represents the notifications table in database
type Notification struct {
	ID             string     `gorm:"type:char(36);primaryKey" json:"id"`
	Type           string     `gorm:"type:varchar(255);not null" json:"type"`
	NotifiableType string     `gorm:"type:varchar(255);not null;index:idx_notifiable" json:"notifiable_type"`
	NotifiableID   uint       `gorm:"not null;index:idx_notifiable" json:"notifiable_id"`
	Data           string     `gorm:"type:text;not null" json:"data"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
