package service

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/baobei23/e-ticket/services/notification-service/internal/domain"
)

//go:embed templates/activation_email.html
var tmplFS embed.FS

type NotificationService struct {
	mailer            domain.Mailer
	activationBaseURL string
	tmpl              *template.Template
	maxRetries        int
	baseDelay         time.Duration
}

func NewNotificationService(mailer domain.Mailer, activationBaseURL string, maxRetries int, baseDelay time.Duration) *NotificationService {
	tmpl := template.Must(template.ParseFS(tmplFS, "templates/activation_email.html"))
	return &NotificationService{
		mailer:            mailer,
		activationBaseURL: activationBaseURL,
		tmpl:              tmpl,
		maxRetries:        maxRetries,
		baseDelay:         baseDelay,
	}
}

func (s *NotificationService) SendActivationEmail(email, token string, expiresAt time.Time) error {
	activationURL := fmt.Sprintf("%s?token=%s", s.activationBaseURL, token)

	data := struct {
		Email         string
		ActivationURL string
		ExpiresAt     string
	}{
		Email:         email,
		ActivationURL: activationURL,
		ExpiresAt:     expiresAt.Format(time.RFC1123),
	}

	var buf bytes.Buffer
	if err := s.tmpl.Execute(&buf, data); err != nil {
		return err
	}

	subject := "Aktivasi Akun Anda"
	return s.sendWithRetry(email, subject, buf.String())
}

func (s *NotificationService) sendWithRetry(to, subject, html string) error {
	var err error
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			delay := s.baseDelay * time.Duration(1<<(attempt-1))
			time.Sleep(delay)
		}

		err = s.mailer.SendHTML(to, subject, html)
		if err == nil {
			return nil
		}
	}
	return err
}
