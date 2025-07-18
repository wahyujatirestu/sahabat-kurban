package service

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	SendVerificationEmail(toEmail, toName, token string) error
	SendResetPasswordEmail(toEmail, toName, token string) error
}

type emailService struct {
	apiKey      string
	fromAddress string
	fromName    string
	appBaseURL  string
}

func NewEmailService(apiKey, fromAddress, fromName, appBaseURL string) EmailService {
	return &emailService{
		apiKey:      apiKey,
		fromAddress: fromAddress,
		fromName:    fromName,
		appBaseURL:  appBaseURL,
	}
}

func (e *emailService) SendVerificationEmail(toEmail, toName, token string) error {
	from := mail.NewEmail(e.fromName, e.fromAddress)
	to := mail.NewEmail(toName, toEmail)

	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", e.appBaseURL, token)

	subject := "Verifikasi Email Sahabat Kurban"
	content := fmt.Sprintf(`
		<h2>Halo %s!</h2>
		<p>Terima kasih telah mendaftar. Silakan klik link berikut untuk verifikasi email Anda:</p>
		<p><a href="%s">%s</a></p>
		<p>Jika Anda tidak merasa membuat akun, abaikan email ini.</p>
	`, toName, verifyURL, verifyURL)

	message := mail.NewSingleEmail(from, subject, to, "", content)

	client := sendgrid.NewSendClient(e.apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println("SendGrid error:", err)
		return err
	}

	log.Printf("Email sent - Status: %d\nBody: %s\n", response.StatusCode, response.Body)
	return nil
}

func (e *emailService) SendResetPasswordEmail(toEmail, toName, token string) error {
	from := mail.NewEmail(e.fromName, e.fromAddress)
	to := mail.NewEmail(toName, toEmail)

	resetURL := fmt.Sprintf("%s/api/v1/auth/reset-password?token=%s", e.appBaseURL, token)
	subject := "Reset Password Sahabat Kurban"

	content := fmt.Sprintf(`
		<h2>Halo %s!</h2>
		<p>Silakan klik link berikut untuk reset password Anda:</p>
		<p><a href="%s">%s</a></p>
		<p>Link hanya berlaku selama 30 menit.</p>
	`, toName, resetURL, resetURL)

	message := mail.NewSingleEmail(from, subject, to, "", content)
	client := sendgrid.NewSendClient(e.apiKey)
	_, err := client.Send(message)
	return err
}


