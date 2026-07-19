package resources

import (
	"encoding/json"

	"finance-app-backend/internal/models"
)

type NotificationResource struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Message   string  `json:"message"`
	Type      string  `json:"type"`
	IsRead    bool    `json:"is_read"`
	Data      string  `json:"data"`
	CreatedAt string  `json:"created_at"`
}

func ToNotificationResource(n *models.Notification) NotificationResource {
	title := "Notifikasi"
	message := ""

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(n.Data), &dataMap); err == nil {
		if t, ok := dataMap["title"].(string); ok {
			title = t
		}
		if m, ok := dataMap["message"].(string); ok {
			message = m
		}
	} else {
		message = n.Data
	}

	return NotificationResource{
		ID:        n.ID,
		Title:     title,
		Message:   message,
		Type:      n.Type,
		IsRead:    n.ReadAt != nil,
		Data:      n.Data,
		CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToNotificationCollection(notes []models.Notification) []NotificationResource {
	res := make([]NotificationResource, len(notes))
	for i, n := range notes {
		res[i] = ToNotificationResource(&n)
	}
	return res
}
