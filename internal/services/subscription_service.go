package services

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"finance-app-backend/internal/models"
	"finance-app-backend/internal/pkg/fcm"
	"finance-app-backend/internal/repositories"
)

type SubscriptionService struct {
	subRepo  repositories.SubscriptionRepositoryInterface
	userRepo repositories.UserRepositoryInterface
	noteRepo repositories.NotificationRepositoryInterface
	xendit   *XenditService
	fcm      *fcm.FCMService
}

func NewSubscriptionService(
	subRepo repositories.SubscriptionRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	noteRepo repositories.NotificationRepositoryInterface,
	xendit *XenditService,
	fcm *fcm.FCMService,
) *SubscriptionService {
	return &SubscriptionService{
		subRepo:  subRepo,
		userRepo: userRepo,
		noteRepo: noteRepo,
		xendit:   xendit,
		fcm:      fcm,
	}
}

func (s *SubscriptionService) GetPlans() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id":            "monthly",
			"name":          "Paket Bulanan",
			"price":         15000,
			"formatted":     "Rp 15.000 / bulan",
			"duration_days": 30,
			"features": []string{
				"Akses tanpa batas ke semua fitur",
				"Tambah dompet kustom tanpa batas",
				"Tambah kategori kustom tanpa batas",
				"Ekspor laporan ke PDF & Excel",
			},
		},
		{
			"id":            "yearly",
			"name":          "Paket Tahunan",
			"price":         150000,
			"formatted":     "Rp 150.000 / tahun",
			"duration_days": 365,
			"discount":      "Hemat 17%",
			"features": []string{
				"Semua fitur Paket Bulanan",
				"Hemat 2 bulan pembayaran",
				"Prioritas dukungan pelanggan",
			},
		},
	}
}

func (s *SubscriptionService) GetStatus(userID uint) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	isPro := user.CheckIsPro()
	var subUntil *string
	if user.SubscriptionUntil != nil {
		formatted := user.SubscriptionUntil.Format("2006-01-02 15:04:05")
		subUntil = &formatted
	}

	return map[string]interface{}{
		"is_pro":             isPro,
		"subscription_until": subUntil,
	}, nil
}

func (s *SubscriptionService) getPlanAmount(plan string) (float64, int, error) {
	if plan == "monthly" {
		return 15000, 30, nil
	} else if plan == "yearly" {
		return 150000, 365, nil
	}
	return 0, 0, errors.New("paket langganan tidak valid")
}

func (s *SubscriptionService) generateExternalID(userID uint, plan string) string {
	return fmt.Sprintf("SUB-%d-%s-%d", userID, plan, time.Now().Unix())
}

func (s *SubscriptionService) generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (s *SubscriptionService) CreateQrisPayment(userID uint, plan string) (*models.SubscriptionLog, error) {
	amount, _, err := s.getPlanAmount(plan)
	if err != nil {
		return nil, err
	}

	extID := s.generateExternalID(userID, plan)
	resp, err := s.xendit.CreateQrisPayment(extID, amount)

	var checkoutUrl *string
	var xenditID *string

	if err == nil {
		if idStr, ok := resp["id"].(string); ok {
			xenditID = &idStr
		}
		if pm, ok := resp["payment_method"].(map[string]interface{}); ok {
			if qr, ok := pm["qr_code"].(map[string]interface{}); ok {
				if qstr, ok := qr["channel_properties"].(map[string]interface{})["qr_string"].(string); ok {
					checkoutUrl = &qstr
				}
			}
		}
	} else {
		log.Printf("Xendit QRIS error: %v. Using fallback simulation.", err)
		dummyQR := "00020101021226670016COM.XENDIT.WWW0118936009140000000000021500000000000000000303UMI51440014ID.CO.QRIS.WWW0215ID10200000000000303UMI5204581253033605802ID5914PersonalFinance6007Jakarta6304A1B2"
		checkoutUrl = &dummyQR
		xenditID = &extID
	}

	pmName := "qris"
	pmChannel := "QRIS"
	subLog := &models.SubscriptionLog{
		UserID:           userID,
		Type:             "qris",
		PlanID:           &plan,
		Amount:           amount,
		Status:           "pending",
		PaymentMethod:    &pmName,
		PaymentChannel:   &pmChannel,
		XenditInvoiceID:  xenditID,
		XenditInvoiceURL: checkoutUrl,
		StartsAt:         time.Now(),
		EndsAt:           time.Now(),
	}

	if err := s.subRepo.Create(subLog); err != nil {
		return nil, err
	}

	return subLog, nil
}

func (s *SubscriptionService) CreateVaPayment(userID uint, plan string, bankCode string) (*models.SubscriptionLog, error) {
	amount, _, err := s.getPlanAmount(plan)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	extID := s.generateExternalID(userID, plan)
	resp, err := s.xendit.CreateVaPayment(extID, amount, bankCode, user.Name)

	var vaNum *string
	var xenditID *string

	if err == nil {
		if idStr, ok := resp["id"].(string); ok {
			xenditID = &idStr
		}
		if pm, ok := resp["payment_method"].(map[string]interface{}); ok {
			if va, ok := pm["virtual_account"].(map[string]interface{}); ok {
				if num, ok := va["channel_properties"].(map[string]interface{})["virtual_account_number"].(string); ok {
					vaNum = &num
				}
			}
		}
	} else {
		log.Printf("Xendit VA error: %v. Using fallback simulation.", err)
		dummyVA := "8808" + fmt.Sprintf("%010d", userID)
		vaNum = &dummyVA
		xenditID = &extID
	}

	pmName := "virtual_account"
	subLog := &models.SubscriptionLog{
		UserID:           userID,
		Type:             "virtual_account",
		PlanID:           &plan,
		Amount:           amount,
		Status:           "pending",
		PaymentMethod:    &pmName,
		PaymentChannel:   &bankCode,
		XenditInvoiceID:  xenditID,
		XenditInvoiceURL: vaNum,
		StartsAt:         time.Now(),
		EndsAt:           time.Now(),
	}

	if err := s.subRepo.Create(subLog); err != nil {
		return nil, err
	}

	return subLog, nil
}

func (s *SubscriptionService) CreateEwalletPayment(userID uint, plan string, channelCode string) (*models.SubscriptionLog, error) {
	amount, _, err := s.getPlanAmount(plan)
	if err != nil {
		return nil, err
	}

	extID := s.generateExternalID(userID, plan)
	resp, err := s.xendit.CreateEwalletPayment(extID, amount, channelCode)

	var checkoutUrl *string
	var xenditID *string

	if err == nil {
		if idStr, ok := resp["id"].(string); ok {
			xenditID = &idStr
		}
		if actions, ok := resp["actions"].([]interface{}); ok && len(actions) > 0 {
			if act, ok := actions[0].(map[string]interface{}); ok {
				if urlStr, ok := act["url"].(string); ok {
					checkoutUrl = &urlStr
				}
			}
		}
	} else {
		log.Printf("Xendit E-Wallet error: %v. Using fallback simulation.", err)
		dummyURL := "https://checkout.xendit.co/web/" + extID
		checkoutUrl = &dummyURL
		xenditID = &extID
	}

	pmName := "ewallet"
	subLog := &models.SubscriptionLog{
		UserID:           userID,
		Type:             "ewallet",
		PlanID:           &plan,
		Amount:           amount,
		Status:           "pending",
		PaymentMethod:    &pmName,
		PaymentChannel:   &channelCode,
		XenditInvoiceID:  xenditID,
		XenditInvoiceURL: checkoutUrl,
		StartsAt:         time.Now(),
		EndsAt:           time.Now(),
	}

	if err := s.subRepo.Create(subLog); err != nil {
		return nil, err
	}

	return subLog, nil
}

func (s *SubscriptionService) CheckPaymentStatus(userID uint, logID uint) (*models.SubscriptionLog, error) {
	sLog, err := s.subRepo.GetByID(logID)
	if err != nil {
		return nil, err
	}

	if sLog.UserID != userID {
		return nil, errors.New("akses ditolak")
	}

	return sLog, nil
}

func (s *SubscriptionService) HandleXenditWebhook(payload map[string]interface{}, callbackToken string) error {
	// Verify Xendit Callback Token
	expectedToken := os.Getenv("XENDIT_CALLBACK_TOKEN")
	if expectedToken != "" && callbackToken != expectedToken {
		return errors.New("unauthorized webhook token")
	}

	// Extract External ID & Status
	extID, _ := payload["external_id"].(string)
	if extID == "" {
		if ref, ok := payload["reference_id"].(string); ok {
			extID = ref
		}
	}

	if extID == "" {
		return errors.New("missing external_id in webhook payload")
	}

	sLog, err := s.subRepo.GetByExternalID(extID)
	if err != nil || sLog == nil {
		log.Printf("Webhook error: SubscriptionLog with external_id %s not found", extID)
		return nil
	}

	if sLog.Status == "paid" {
		log.Printf("SubscriptionLog %s already marked paid. Skipping.", extID)
		return nil
	}

	statusStr, _ := payload["status"].(string)

	if statusStr == "PAID" || statusStr == "SUCCEEDED" || statusStr == "COMPLETED" {
		return s.ActivateSubscription(sLog)
	} else if statusStr == "EXPIRED" {
		sLog.Status = "expired"
		return s.subRepo.Update(sLog)
	} else if statusStr == "FAILED" {
		sLog.Status = "failed"
		return s.subRepo.Update(sLog)
	}

	return nil
}

func (s *SubscriptionService) ActivateSubscription(sLog *models.SubscriptionLog) error {
	user, err := s.userRepo.GetByID(sLog.UserID)
	if err != nil {
		return err
	}

	plan := "monthly"
	if sLog.PlanID != nil {
		plan = *sLog.PlanID
	}

	_, durationDays, err := s.getPlanAmount(plan)
	if err != nil {
		durationDays = 30
	}

	now := time.Now()
	var startsAt time.Time

	// Stacking duration logic: if currently active Pro, start from current subscription_until
	if user.SubscriptionUntil != nil && user.SubscriptionUntil.After(now) {
		startsAt = *user.SubscriptionUntil
	} else {
		startsAt = now
	}

	endsAt := startsAt.AddDate(0, 0, durationDays)

	// Update user status
	user.IsPro = true
	user.SubscriptionUntil = &endsAt
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Update log status
	sLog.Status = "paid"
	sLog.StartsAt = startsAt
	sLog.EndsAt = endsAt
	if err := s.subRepo.Update(sLog); err != nil {
		return err
	}

	// Create internal notification record
	title := "Pembayaran Berhasil!"
	msg := fmt.Sprintf("Selamat, akun Anda telah diperbarui menjadi Pro hingga %s.", endsAt.Format("02 Jan 2006"))

	dataBytes, _ := json.Marshal(map[string]string{
		"title":   title,
		"message": msg,
	})

	note := &models.Notification{
		ID:             s.generateUUID(),
		Type:           "App\\Notifications\\PaymentSuccessNotification",
		NotifiableType: "App\\Models\\User",
		NotifiableID:   user.ID,
		Data:           string(dataBytes),
		ReadAt:         nil,
	}
	s.noteRepo.Create(note)

	// Send Push Notification if FCM token available
	if user.FcmToken != nil && *user.FcmToken != "" {
		s.fcm.SendToDevice(*user.FcmToken, title, msg, map[string]string{
			"type":            "subscription",
			"subscription_id": fmt.Sprintf("%d", sLog.ID),
		})
	}

	return nil
}

func (s *SubscriptionService) GetHistory(userID uint) ([]models.SubscriptionLog, error) {
	return s.subRepo.GetForUser(userID)
}
