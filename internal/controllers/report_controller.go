package controllers

import (
	"strconv"
	"time"

	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ReportController struct {
	reportService *services.ReportService
}

func NewReportController(reportService *services.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (ctrl *ReportController) Monthly(c *fiber.Ctx) error {
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

	data, err := ctrl.reportService.GetMonthlyReport(userID, month, year)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

func (ctrl *ReportController) Category(c *fiber.Ctx) error {
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

	data, err := ctrl.reportService.GetCategoryReport(userID, month, year)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(data)
}
