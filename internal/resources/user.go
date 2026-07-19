package resources

import (
	"fmt"
	"strings"

	"finance-app-backend/internal/models"
)

// UserResponse defines user payload format returned to the mobile app
type UserResponse struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Email             string  `json:"email"`
	PhotoURL          *string `json:"photo_url"`
	Phone             *string `json:"phone"`
	Address           *string `json:"address"`
	DateOfBirth       *string `json:"date_of_birth"`
	Gender            *string `json:"gender"`
	IsPro             bool    `json:"is_pro"`
	SubscriptionUntil *string `json:"subscription_until"`
}

// ToUserResponse maps a GORM User model to UserResponse
func ToUserResponse(user *models.User, appURL string) UserResponse {
	var dob *string
	if user.DateOfBirth != nil {
		dateStr := user.DateOfBirth.Format("2006-01-02")
		dob = &dateStr
	}

	var subUntil *string
	if user.SubscriptionUntil != nil {
		subStr := user.SubscriptionUntil.Format("2006-01-02T15:04:05.000000Z")
		subUntil = &subStr
	}

	// Resolve photo URL (local filesystem or full cloud url)
	var resolvedPhotoURL *string
	if user.PhotoURL != nil {
		photo := *user.PhotoURL
		if strings.HasPrefix(photo, "http") {
			resolvedPhotoURL = &photo
		} else {
			fullURL := fmt.Sprintf("%s/storage/%s", appURL, strings.TrimPrefix(photo, "/"))
			resolvedPhotoURL = &fullURL
		}
	}

	return UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		PhotoURL:          resolvedPhotoURL,
		Phone:             user.Phone,
		Address:           user.Address,
		DateOfBirth:       dob,
		Gender:            user.Gender,
		IsPro:             user.CheckIsPro(),
		SubscriptionUntil: subUntil,
	}
}
