package controllers

import (
	"log"

	"finance-app-backend/internal/config"
	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/pkg/validator"
	"finance-app-backend/internal/requests"
	"finance-app-backend/internal/resources"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	user, token, err := ctrl.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToUserResponse(user, config.AppConfig.AppUrl)
	return response.Success(c, "Pendaftaran berhasil", fiber.Map{
		"user":  res,
		"token": token,
	})
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	user, token, err := ctrl.authService.Login(req.Email, req.Password)
	if err != nil {
		// Matching Laravel exact error response shape
		errs := map[string][]string{"email": {err.Error()}}
		return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), errs)
	}

	res := resources.ToUserResponse(user, config.AppConfig.AppUrl)
	return response.Success(c, "Login berhasil", fiber.Map{
		"user":  res,
		"token": token,
	})
}

func (ctrl *AuthController) SocialLogin(c *fiber.Ctx) error {
	var req requests.SocialLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	user, token, err := ctrl.authService.SocialLogin(req.Provider, req.IdToken, req.Name, req.Email)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	res := resources.ToUserResponse(user, config.AppConfig.AppUrl)
	return response.Success(c, "Login sosial berhasil", fiber.Map{
		"user":  res,
		"token": token,
	})
}

func (ctrl *AuthController) SendOtp(c *fiber.Ctx) error {
	var req requests.SendOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	err := ctrl.authService.SendOtp(req.Email)
	if err != nil {
		errs := map[string][]string{"email": {err.Error()}}
		return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), errs)
	}

	return response.SuccessNoContent(c, "Kode OTP telah dikirim ke email Anda.")
}

func (ctrl *AuthController) VerifyOtp(c *fiber.Ctx) error {
	var req requests.VerifyOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	err := ctrl.authService.VerifyOtp(req.Email, req.Otp)
	if err != nil {
		errs := map[string][]string{"otp": {err.Error()}}
		return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), errs)
	}

	return response.SuccessNoContent(c, "Kode OTP valid.")
}

func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	var req requests.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	err := ctrl.authService.ResetPassword(req.Email, req.Otp, req.Password)
	if err != nil {
		errs := map[string][]string{"otp": {err.Error()}}
		return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), errs)
	}

	return response.SuccessNoContent(c, "Password berhasil diubah. Silakan login dengan password baru.")
}

func (ctrl *AuthController) GetUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	user, err := ctrl.authService.GetUserByID(userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Pengguna tidak ditemukan", nil)
	}

	res := resources.ToUserResponse(user, config.AppConfig.AppUrl)
	return response.Success(c, "Data profile berhasil diambil", res)
}

func (ctrl *AuthController) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req requests.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	user, err := ctrl.authService.UpdateProfile(userID, req.Name, req.Phone, req.Address, req.DateOfBirth, req.Gender)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToUserResponse(user, config.AppConfig.AppUrl)
	return response.Success(c, "Profil berhasil diperbarui", res)
}

func (ctrl *AuthController) UploadPhoto(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	file, err := c.FormFile("photo")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "File foto tidak valid atau tidak diterima.", nil)
	}

	// Limit photo upload to 2MB
	if file.Size > 2048*1024 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Ukuran file maksimal adalah 2MB.", nil)
	}

	fileHeader, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Gagal memproses file foto.", nil)
	}
	defer fileHeader.Close()

	user, err := ctrl.authService.SavePhoto(userID, fileHeader, file.Filename)
	if err != nil {
		log.Printf("UploadPhoto error: %v", err)
		return response.Error(c, fiber.StatusInternalServerError, "Gagal menyimpan foto. Silakan coba lagi.", nil)
	}

	res := resources.ToUserResponse(user, config.AppConfig.AppUrl)
	return response.Success(c, "Foto profil berhasil diperbarui", res)
}

func (ctrl *AuthController) DeletePhoto(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	user, err := ctrl.authService.DeletePhoto(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToUserResponse(user, config.AppConfig.AppUrl)
	return response.Success(c, "Foto profil berhasil dihapus", res)
}

func (ctrl *AuthController) UpdateFcmToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req requests.FcmTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	err := ctrl.authService.UpdateFcmToken(userID, req.FcmToken)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.SuccessNoContent(c, "FCM token berhasil diperbarui")
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	err := ctrl.authService.ClearFcmTokenAndLogout(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.SuccessNoContent(c, "Logged out")
}
