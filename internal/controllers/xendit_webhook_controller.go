package controllers

import (
	"log"

	"finance-app-backend/internal/pkg/response"
	"finance-app-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type XenditWebhookController struct {
	subService *services.SubscriptionService
}

func NewXenditWebhookController(subService *services.SubscriptionService) *XenditWebhookController {
	return &XenditWebhookController{subService: subService}
}

func (ctrl *XenditWebhookController) HandleWebhook(c *fiber.Ctx) error {
	callbackToken := c.Get("x-callback-token")

	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Xendit Webhook payload parse error: %v", err)
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload", nil)
	}

	err := ctrl.subService.HandleXenditWebhook(payload, callbackToken)
	if err != nil {
		log.Printf("Xendit Webhook error: %v", err)
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, "Webhook processed successfully", fiber.Map{"status": "ok"})
}
