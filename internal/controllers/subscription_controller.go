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

type SubscriptionController struct {
	subService *services.SubscriptionService
}

func NewSubscriptionController(subService *services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{subService: subService}
}

func (ctrl *SubscriptionController) Plans(c *fiber.Ctx) error {
	plans := ctrl.subService.GetPlans()
	return response.Success(c, "Daftar paket langganan", plans)
}

func (ctrl *SubscriptionController) Status(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	status, err := ctrl.subService.GetStatus(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, "Status langganan", status)
}

func (ctrl *SubscriptionController) PayQris(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req requests.QrisPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	sLog, err := ctrl.subService.CreateQrisPayment(userID, req.Plan)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToSubscriptionLogResource(sLog)
	return response.Success(c, "Pembayaran QRIS berhasil dibuat", res)
}

func (ctrl *SubscriptionController) PayVa(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req requests.VaPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	sLog, err := ctrl.subService.CreateVaPayment(userID, req.Plan, req.BankCode)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToSubscriptionLogResource(sLog)
	return response.Success(c, "Pembayaran VA berhasil dibuat", res)
}

func (ctrl *SubscriptionController) PayEwallet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req requests.EwalletPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid payload format", nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "Validasi gagal", errs)
	}

	sLog, err := ctrl.subService.CreateEwalletPayment(userID, req.Plan, req.ChannelCode)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToSubscriptionLogResource(sLog)
	return response.Success(c, "Pembayaran E-Wallet berhasil dibuat", res)
}

func (ctrl *SubscriptionController) CheckStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	sLog, err := ctrl.subService.CheckPaymentStatus(userID, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
	}

	res := resources.ToSubscriptionLogResource(sLog)
	return response.Success(c, "Status pembayaran ditemukan", res)
}

func (ctrl *SubscriptionController) History(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	logs, err := ctrl.subService.GetHistory(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToSubscriptionLogCollection(logs)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
