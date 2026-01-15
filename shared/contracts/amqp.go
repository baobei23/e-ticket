package contracts

import (
	"encoding/json"
	"time"
)

// Nama Queue
const (
	QueueBookingCreated          = "booking.created"
	QueuePaymentSuccess          = "payment.success"
	QueuePaymentFailed           = "payment.failed"
	QueueUserActivationRequested = "user.activation_requested"
)

// Wrapper Message Standar
type AmqpMessage struct {
	EventName string          `json:"event_name"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// Payload: Booking Created
type BookingCreatedEvent struct {
	BookingID string    `json:"booking_id"`
	UserID    int64     `json:"user_id"`
	EventID   int64     `json:"event_id"`
	Quantity  int32     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

// Payload: Payment Success
type PaymentSuccessEvent struct {
	BookingID string  `json:"booking_id"`
	PaymentID string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
}

// Payload: Payment Failed
type PaymentFailedEvent struct {
	BookingID string `json:"booking_id"`
	Reason    string `json:"reason"`
}

// Payload: User Activation Requested
type UserActivationRequestedEvent struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
