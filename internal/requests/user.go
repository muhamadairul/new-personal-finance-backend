package requests

// UpdateProfileRequest holds user profile update payload
type UpdateProfileRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Phone       *string `json:"phone" validate:"omitempty,max=20"`
	Address     *string `json:"address" validate:"omitempty,max=500"`
	DateOfBirth *string `json:"date_of_birth" validate:"omitempty"` // Expected format YYYY-MM-DD
	Gender      *string `json:"gender" validate:"omitempty,oneof=L P"`
}

// FcmTokenRequest holds registration token for push notification
type FcmTokenRequest struct {
	FcmToken string `json:"fcm_token" validate:"required,max=255"`
}
