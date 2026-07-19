package controllers

import (
	"strconv"

	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/resources"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type NotificationController struct {
	noteService *services.NotificationService
}

func NewNotificationController(noteService *services.NotificationService) *NotificationController {
	return &NotificationController{noteService: noteService}
}

func (ctrl *NotificationController) Index(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	notes, total, err := ctrl.noteService.List(userID, page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	res := resources.ToNotificationCollection(notes)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
		},
	})
}

func (ctrl *NotificationController) UnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	count, err := ctrl.noteService.CountUnread(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, "Jumlah notifikasi belum dibaca", fiber.Map{
		"unread_count": count,
	})
}

func (ctrl *NotificationController) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "ID tidak valid", nil)
	}

	err := ctrl.noteService.MarkRead(userID, id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.SuccessNoContent(c, "Notifikasi telah ditandai dibaca.")
}

func (ctrl *NotificationController) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	err := ctrl.noteService.MarkAllRead(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.SuccessNoContent(c, "Semua notifikasi telah ditandai dibaca.")
}
