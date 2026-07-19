package controllers

import (
	"strconv"

	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/pkg/validator"
	"finance-app-backend/internal/requests"
	"finance-app-backend/internal/resources"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type WalletController struct {
	walletService *services.WalletService
}

func NewWalletController(walletService *services.WalletService) *WalletController {
	return &WalletController{walletService: walletService}
}

func (ctrl *WalletController) Index(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	wallets, err := ctrl.walletService.List(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToWalletCollection(wallets)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (ctrl *WalletController) Store(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req requests.WalletRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	balance := 0.0
	if req.Balance != nil {
		balance = *req.Balance
	}

	wallet, err := ctrl.walletService.Create(userID, req.Name, req.Type, balance, req.Icon, req.Color)
	if err != nil {
		// Matching Laravel 403 response for Free users trying to add wallet limit
		return response.Error(c, fiber.StatusForbidden, err.Error(), fiber.Map{
			"upgrade_required": true,
		})
	}

	res := resources.ToWalletResource(wallet)
	return response.Success(c, "Dompet berhasil dibuat", res)
}

func (ctrl *WalletController) Show(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	wallet, err := ctrl.walletService.GetByID(userID, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, "Akses ditolak", nil)
	}

	res := resources.ToWalletResource(wallet)
	return response.Success(c, "Dompet ditemukan", res)
}

func (ctrl *WalletController) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	var req requests.WalletRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	balance := 0.0
	if req.Balance != nil {
		balance = *req.Balance
	}

	wallet, err := ctrl.walletService.Update(userID, uint(id), req.Name, req.Type, balance, req.Icon, req.Color)
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
	}

	res := resources.ToWalletResource(wallet)
	return response.Success(c, "Dompet berhasil diperbarui", res)
}

func (ctrl *WalletController) Destroy(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	err = ctrl.walletService.Delete(userID, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
	}

	return response.SuccessNoContent(c, "Dompet berhasil dihapus")
}
