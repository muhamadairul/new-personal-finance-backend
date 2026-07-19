package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type FCMService struct {
	projectID       string
	credentialsPath string
}

func NewFCMService() *FCMService {
	return &FCMService{
		projectID:       os.Getenv("FIREBASE_PROJECT_ID"),
		credentialsPath: os.Getenv("FIREBASE_CREDENTIALS"),
	}
}

// SendToDevice sends a push notification to a single FCM device token
func (s *FCMService) SendToDevice(fcmToken, title, body string, extraData map[string]string) bool {
	if fcmToken == "" {
		return false
	}

	if s.projectID == "" {
		log.Printf("[FCM LOG] Mock push to token '%s...': Title='%s', Body='%s'", fcmToken[:min(20, len(fcmToken))], title, body)
		return true
	}

	// Prepare payload
	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"token": fcmToken,
			"notification": map[string]string{
				"title": title,
				"body":  body,
			},
			"android": map[string]interface{}{
				"priority": "high",
				"notification": map[string]interface{}{
					"channel_id":              "finance_notifications",
					"default_sound":           true,
					"default_vibrate_timings": true,
				},
			},
			"data": extraData,
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("FCM marshal error: %v", err)
		return false
	}

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", s.projectID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("FCM request error: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	// If access token is missing in dev environment, log gracefully
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("FCM send error: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
