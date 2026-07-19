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

type BudgetController struct {
	budgetService *services.BudgetService
}

func NewBudgetController(budgetService *services.BudgetService) *BudgetController {
	return &BudgetController{budgetService: budgetService}
}

func (ctrl *BudgetController) Index(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	now := time.Now()
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	month := int(now.Month())
	year := now.Year()

	if monthStr != "" {
		month, _ = strconv.Atoi(monthStr)
	}
	if yearStr != "" {
		year, _ = strconv.Atoi(yearStr)
	}

	budgets, err := ctrl.budgetService.List(userID, month, year)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToBudgetCollection(budgets)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (ctrl *BudgetController) Store(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req requests.BudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	budget, err := ctrl.budgetService.Upsert(userID, req.CategoryID, req.Amount, req.Month, req.Year)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToBudgetResource(budget)
	return response.Success(c, "Anggaran berhasil disimpan", res)
}

func (ctrl *BudgetController) Show(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	budget, err := ctrl.budgetService.GetByID(userID, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, "Akses ditolak", nil)
	}

	res := resources.ToBudgetResource(budget)
	return response.Success(c, "Anggaran ditemukan", res)
}

func (ctrl *BudgetController) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	var req requests.BudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	budget, err := ctrl.budgetService.Update(userID, uint(id), req.CategoryID, req.Amount, req.Month, req.Year)
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
	}

	res := resources.ToBudgetResource(budget)
	return response.Success(c, "Anggaran berhasil diperbarui", res)
}

func (ctrl *BudgetController) Destroy(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	err = ctrl.budgetService.Delete(userID, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
	}

	return response.SuccessNoContent(c, "Anggaran berhasil dihapus")
}
