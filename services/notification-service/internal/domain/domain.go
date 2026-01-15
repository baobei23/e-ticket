package domain

import "time"

type Mailer interface {
	SendHTML(to, subject, html string) error
}

type NotificationService interface {
	SendActivationEmail(email, token string, expiresAt time.Time) error
}
