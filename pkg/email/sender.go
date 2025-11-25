package email

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/taiwanstay/taiwanstay-back/pkg/config"
	"github.com/taiwanstay/taiwanstay-back/pkg/logger"
)

// EmailSender 定義發送郵件的介面
type EmailSender interface {
	Send(toEmail, toName, subject, htmlBody string) error
}

// BrevoSender 實作 Brevo (Sendinblue) 郵件發送
type BrevoSender struct {
	apiKey      string
	senderEmail string
	senderName  string
	client      *http.Client
}

func NewBrevoSender(cfg *config.Config) *BrevoSender {
	return &BrevoSender{
		apiKey:      cfg.Email.BrevoAPIKey,
		senderEmail: cfg.Email.BrevoSenderEmail,
		senderName:  cfg.Email.BrevoSenderName,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *BrevoSender) Send(toEmail, toName, subject, htmlBody string) error {
	url := "https://api.brevo.com/v3/smtp/email"

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  s.senderName,
			"email": s.senderEmail,
		},
		"to": []map[string]string{
			{
				"email": toEmail,
				"name":  toName,
			},
		},
		"subject":     subject,
		"htmlContent": htmlBody,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", s.apiKey)
	req.Header.Set("content-type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("brevo api error: status code %d", resp.StatusCode)
	}

	return nil
}

// MailerLiteSender 實作 MailerLite 郵件發送 (Fallback)
type MailerLiteSender struct {
	apiKey      string
	senderEmail string // MailerLite requires verified sender
	senderName  string
	client      *http.Client
}

func NewMailerLiteSender(cfg *config.Config) *MailerLiteSender {
	return &MailerLiteSender{
		apiKey:      cfg.Email.MailerLiteAPIKey,
		senderEmail: cfg.Email.BrevoSenderEmail, // Assuming same sender email is verified on both
		senderName:  cfg.Email.BrevoSenderName,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *MailerLiteSender) Send(toEmail, toName, subject, htmlBody string) error {
	// MailerLite Transactional API (New)
	url := "https://connect.mailerlite.com/api/emails"

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": s.senderEmail,
			"name":  s.senderName,
		},
		"to": []map[string]string{
			{
				"email": toEmail,
				"name":  toName,
			},
		},
		"subject": subject,
		"html":    htmlBody,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("mailerlite api error: status code %d", resp.StatusCode)
	}

	return nil
}

// FallbackSender 實作備援機制
type FallbackSender struct {
	primary   EmailSender
	secondary EmailSender
}

func NewFallbackSender(primary, secondary EmailSender) *FallbackSender {
	return &FallbackSender{
		primary:   primary,
		secondary: secondary,
	}
}

func (s *FallbackSender) Send(toEmail, toName, subject, htmlBody string) error {
	// Try Primary
	err := s.primary.Send(toEmail, toName, subject, htmlBody)
	if err == nil {
		logger.Info("Email sent via Primary (Brevo)", "to", toEmail)
		return nil
	}

	logger.Warn("Primary email sender failed, trying secondary", "error", err)

	// Try Secondary
	err = s.secondary.Send(toEmail, toName, subject, htmlBody)
	if err == nil {
		logger.Info("Email sent via Secondary (MailerLite)", "to", toEmail)
		return nil
	}

	logger.Error("Both email senders failed", "error", err)
	return errors.New("failed to send email via both providers")
}
