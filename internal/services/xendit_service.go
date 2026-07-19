package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"finance-app-backend/internal/config"
)

type XenditService struct {
	secretKey string
	client    *http.Client
}

func NewXenditService() *XenditService {
	return &XenditService{
		secretKey: config.AppConfig.JwtSecret, // Fallback if env XENDIT_SECRET_KEY empty
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *XenditService) getSecretKey() string {
	key := os.Getenv("XENDIT_SECRET_KEY")
	if key != "" {
		return key
	}
	return s.secretKey
}

func (s *XenditService) CreateQrisPayment(referenceID string, amount float64) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"reference_id": referenceID,
		"currency":     "IDR",
		"country":      "ID",
		"amount":       amount,
		"payment_method": map[string]interface{}{
			"reference_id": referenceID,
			"type":         "QR_CODE",
			"reusability":  "ONE_TIME_USE",
			"qr_code": map[string]interface{}{
				"channel_code": "QRIS",
				"channel_properties": map[string]interface{}{
					"qr_code_generator": "INTEGRATED",
				},
			},
		},
		"checkout_method": "ONE_TIME_PAYMENT",
	}

	return s.sendPaymentRequest(payload)
}

func (s *XenditService) CreateVaPayment(referenceID string, amount float64, bankCode string, customerName string) (map[string]interface{}, error) {
	expiresAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	payload := map[string]interface{}{
		"reference_id": referenceID,
		"currency":     "IDR",
		"country":      "ID",
		"amount":       amount,
		"payment_method": map[string]interface{}{
			"type":         "VIRTUAL_ACCOUNT",
			"reusability":  "ONE_TIME_USE",
			"reference_id": referenceID,
			"virtual_account": map[string]interface{}{
				"channel_code": bankCode,
				"channel_properties": map[string]interface{}{
					"customer_name": customerName,
					"expires_at":    expiresAt,
				},
			},
		},
		"checkout_method": "ONE_TIME_PAYMENT",
		"metadata": map[string]interface{}{
			"source": "subscription",
		},
	}

	return s.sendPaymentRequest(payload)
}

func (s *XenditService) CreateEwalletPayment(referenceID string, amount float64, channelCode string) (map[string]interface{}, error) {
	appURL := config.AppConfig.AppUrl

	payload := map[string]interface{}{
		"reference_id": referenceID,
		"currency":     "IDR",
		"country":      "ID",
		"amount":       amount,
		"payment_method": map[string]interface{}{
			"type":         "EWALLET",
			"reusability":  "ONE_TIME_USE",
			"reference_id": referenceID,
			"ewallet": map[string]interface{}{
				"channel_code": channelCode,
				"channel_properties": map[string]interface{}{
					"success_return_url": appURL + "/payment/success",
					"failure_return_url": appURL + "/payment/failed",
				},
			},
		},
		"checkout_method": "ONE_TIME_PAYMENT",
	}

	return s.sendPaymentRequest(payload)
}

func (s *XenditService) sendPaymentRequest(payload map[string]interface{}) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.xendit.co/payment_requests", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(s.getSecretKey(), "")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		msg, _ := result["message"].(string)
		if msg == "" {
			msg = fmt.Sprintf("Xendit API returned status %d", resp.StatusCode)
		}
		return nil, errors.New(msg)
	}

	return result, nil
}
