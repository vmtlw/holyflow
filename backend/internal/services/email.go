package services

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/holyflow/backend/internal/config"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		cfg: cfg,
	}
}

func (s *EmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	// If SMTP is not configured, just print to console
	if s.cfg.SMTPHost == "" || s.cfg.SMTPUsername == "" || s.cfg.SMTPPassword == "" {
		fmt.Printf("Mock email sent:\nTo: %s\nSubject: %s\nBody: %s\n", to, subject, body)
		return nil
	}

	// Create message
	message := fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		to, subject, body)

	// Connect to SMTP server
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	// Send email
	err := smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
