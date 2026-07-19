package controllers

import (
	"strconv"
	"time"

	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/pkg/validator"
	"finance-app-backend/internal/requests"
	"finance-app-backend/internal/resources"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	txService *services.TransactionService
}

func NewTransactionController(txService *services.TransactionService) *TransactionController {
	return &TransactionController{txService: txService}
}

func (ctrl *TransactionController) Index(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	monthStr := c.Query("month")
	yearStr := c.Query("year")

	month := 0
	year := 0
	if monthStr != "" && yearStr != "" {
		month, _ = strconv.Atoi(monthStr)
		year, _ = strconv.Atoi(yearStr)
	}

	txs, err := ctrl.txService.List(userID, month, year)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToTransactionCollection(txs)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (ctrl *TransactionController) Store(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req requests.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	date, err := time.Parse("2006-01-02 15:04:05", req.Date)
	if err != nil {
		// Fallback parse common standard formats if not exact YYYY-MM-DD HH:MM:SS
		date, err = time.Parse(time.RFC3339, req.Date)
		if err != nil {
			date, err = time.Parse("2006-01-02", req.Date)
			if err != nil {
				return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", map[string][]string{
					"date": {"Format tanggal tidak valid. Gunakan format YYYY-MM-DD HH:mm:ss"},
				})
			}
		}
	}

	t, err := ctrl.txService.Create(userID, req.Type, req.Amount, req.CategoryID, req.WalletID, req.Note, date)
	if err != nil {
		// Return 422 if it's the balance validation failure
		errs := map[string][]string{"amount": {err.Error()}}
		return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), errs)
	}

	res := resources.ToTransactionResource(t)
	return response.Success(c, "Transaksi berhasil dibuat", res)
}

func (ctrl *TransactionController) Show(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	t, err := ctrl.txService.GetByID(userID, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, "Akses ditolak", nil)
	}

	res := resources.ToTransactionResource(t)
	return response.Success(c, "Transaksi ditemukan", res)
}

func (ctrl *TransactionController) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	var req requests.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	date, err := time.Parse("2006-01-02 15:04:05", req.Date)
	if err != nil {
		date, err = time.Parse(time.RFC3339, req.Date)
		if err != nil {
			date, err = time.Parse("2006-01-02", req.Date)
			if err != nil {
				return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", map[string][]string{
					"date": {"Format tanggal tidak valid."},
				})
			}
		}
	}

	t, err := ctrl.txService.Update(userID, uint(id), req.Type, req.Amount, req.CategoryID, req.WalletID, req.Note, date)
	if err != nil {
		errs := map[string][]string{"amount": {err.Error()}}
		return response.Error(c, fiber.StatusUnprocessableEntity, err.Error(), errs)
	}

	res := resources.ToTransactionResource(t)
	return response.Success(c, "Transaksi berhasil diperbarui", res)
}

func (ctrl *TransactionController) Destroy(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	err = ctrl.txService.Delete(userID, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
	}

	return response.SuccessNoContent(c, "Transaksi berhasil dihapus")
}
