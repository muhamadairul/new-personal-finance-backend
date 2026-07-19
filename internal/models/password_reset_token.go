package models

import (
	"time"
)

// PasswordResetToken represents the password_reset_tokens table in database
type PasswordResetToken struct {
	Email     string     `gorm:"type:varchar(255);primaryKey" json:"email"`
	Token     string     `gorm:"type:varchar(255);not null" json:"token"`
	CreatedAt *time.Time `json:"created_at"`
}
