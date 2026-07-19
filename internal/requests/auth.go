package requests

// RegisterRequest holds registration data
type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,max=255"`
	Email                string `json:"email" validate:"required,email,max=255"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

// LoginRequest holds login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SocialLoginRequest holds Google token login data
type SocialLoginRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google"`
	IdToken  string `json:"id_token" validate:"required"`
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

// SendOtpRequest holds forgot password request email
type SendOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyOtpRequest holds OTP verification data
type VerifyOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,len=6"`
}

// ResetPasswordRequest holds password reset payload
type ResetPasswordRequest struct {
	Email                string `json:"email" validate:"required,email"`
	Otp                  string `json:"otp" validate:"required,len=6"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}
