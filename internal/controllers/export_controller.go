package controllers

import (
	"strconv"
	"time"

	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ExportController struct {
	exportService *services.ExportService
	authService   *services.AuthService
}

func NewExportController(exportService *services.ExportService, authService *services.AuthService) *ExportController {
	return &ExportController{
		exportService: exportService,
		authService:   authService,
	}
}

func (ctrl *ExportController) Excel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Check Pro status
	user, err := ctrl.authService.GetUserByID(userID)
	if err != nil || !user.CheckIsPro() {
		return response.Error(c, fiber.StatusForbidden, "Fitur ini hanya untuk pengguna Pro.", fiber.Map{
			"upgrade_required": true,
		})
	}

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

	fileBytes, filename, err := ctrl.exportService.GenerateExcel(userID, month, year)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengekspor Excel: "+err.Error(), nil)
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	return c.Send(fileBytes)
}

func (ctrl *ExportController) PDF(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Check Pro status
	user, err := ctrl.authService.GetUserByID(userID)
	if err != nil || !user.CheckIsPro() {
		return response.Error(c, fiber.StatusForbidden, "Fitur ini hanya untuk pengguna Pro.", fiber.Map{
			"upgrade_required": true,
		})
	}

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

	fileBytes, filename, err := ctrl.exportService.GeneratePDF(userID, month, year)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal mengekspor PDF: "+err.Error(), nil)
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	return c.Send(fileBytes)
}
