package domain

import "context"

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
	PaymentStatusFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID         string
	BookingID  string
	UserID     int64
	Amount     float64
	Currency   string
	Status     PaymentStatus
	StripeID   string
	PaymentURL string
}

// Interface to External Payment Gateway (Stripe)
type PaymentGateway interface {
	CreateSession(ctx context.Context, payment *Payment) (string, string, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	UpdateStatus(ctx context.Context, id string, status PaymentStatus) error
	GetByBookingID(ctx context.Context, bookingID string) (*Payment, error)
}

type PaymentService interface {
	CreatePayment(ctx context.Context, bookingID string, userID int64, amount float64) (*Payment, error)
	HandleWebhook(ctx context.Context, payload []byte, sigHeader string) error
}
