package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

// SendOtpEmail sends password reset OTP code to user's email
func SendOtpEmail(toEmail, otp string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_ADDRESS")
	fromName := os.Getenv("SMTP_FROM_NAME")

	if fromName == "" {
		fromName = "Pencatat Keuangan"
	}

	// Dynamic year for footer
	currentYear := "2026"

	// HTML email template
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Reset Password</title>
    <style>
        body { font-family: sans-serif; background-color: #f4f4f4; padding: 20px; }
        .container { background-color: #ffffff; padding: 30px; border-radius: 8px; max-width: 500px; margin: 0 auto; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 20px; }
        .otp-box { background-color: #f8fafc; border: 2px dashed #cbd5e1; padding: 15px; text-align: center; font-size: 24px; font-weight: bold; letter-spacing: 4px; color: #0f172a; margin: 20px 0; border-radius: 6px; }
        .footer { text-align: center; margin-top: 30px; font-size: 12px; color: #64748b; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Permintaan Reset Password</h2>
        </div>
        <p>Halo,</p>
        <p>Kami menerima permintaan untuk mereset password akun Anda. Silakan gunakan kode OTP di bawah ini untuk melanjutkan proses reset password. Kode ini hanya berlaku selama 15 menit.</p>
        
        <div class="otp-box">
            {{.OTP}}
        </div>
        
        <p>Jika Anda tidak meminta reset password, Anda dapat mengabaikan email ini dengan aman. Password Anda tidak akan diubah.</p>
        
        <div class="footer">
            &copy; {{.Year}} {{.FromName}}. All rights reserved.
        </div>
    </div>
</body>
</html>`

	// Parse template
	t, err := template.New("otp").Parse(tmpl)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = t.Execute(&body, map[string]string{
		"OTP":      otp,
		"Year":     currentYear,
		"FromName": fromName,
	})
	if err != nil {
		return err
	}

	// If SMTP parameters are missing, fallback to logging
	if smtpHost == "" || smtpUser == "" {
		log.Printf("[MAIL LOG] Fallback send to %s. OTP Code: %s", toEmail, otp)
		return nil
	}

	// Construct message headers
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Kode Verifikasi Reset Password\n"
	fromHeader := fmt.Sprintf("From: %s <%s>\n", fromName, fromEmail)
	toHeader := fmt.Sprintf("To: %s\n", toEmail)

	msg := []byte(fromHeader + toHeader + subject + mime + body.String())

	// Authenticate
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send SMTP email: %w", err)
	}

	return nil
}
