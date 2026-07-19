package models

import (
	"time"
)

// ApiLog represents the api_logs table in database
type ApiLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     *uint     `gorm:"index" json:"user_id,omitempty"`
	Method     string    `gorm:"type:varchar(10);not null" json:"method"`
	URL        string    `gorm:"type:varchar(500);not null" json:"url"`
	Payload    *string   `gorm:"type:text" json:"payload,omitempty"`
	Response   *string   `gorm:"type:text" json:"response,omitempty"`
	StatusCode int       `gorm:"not null" json:"status_code"`
	IPAddress  *string   `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent  *string   `gorm:"type:text" json:"user_agent,omitempty"`
	DurationMs *float64  `gorm:"type:decimal(10,2)" json:"duration_ms,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
