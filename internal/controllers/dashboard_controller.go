package controllers

import (
	"finance-app-backend/internal/config"
	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardController struct {
	reportService *services.ReportService
}

func NewDashboardController(reportService *services.ReportService) *DashboardController {
	return &DashboardController{reportService: reportService}
}

func (ctrl *DashboardController) Index(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	data, err := ctrl.reportService.GetDashboardData(userID, config.AppConfig.AppUrl)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(data)
}
