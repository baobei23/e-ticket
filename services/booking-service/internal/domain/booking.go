package domain

import (
	"context"
	"time"

	pb "github.com/baobei23/e-ticket/shared/proto/booking"
)

// Enum Status
const (
	StatusPending   = "PENDING"
	StatusConfirmed = "CONFIRMED"
	StatusFailed    = "FAILED"
	StatusCancelled = "CANCELLED"
)

type Booking struct {
	ID          string // UUID
	UserID      int64
	EventID     int64
	Quantity    int32
	TotalAmount float64
	Status      string
	CreatedAt   time.Time
}

func (b *Booking) ToProto() *pb.Booking {
	return &pb.Booking{
		BookingId:   b.ID,
		UserId:      b.UserID,
		EventId:     b.EventID,
		Quantity:    b.Quantity,
		TotalAmount: b.TotalAmount,
		Status:      b.Status,
		CreatedAt:   b.CreatedAt.Unix(),
	}
}

// Interfaces

type BookingRepository interface {
	Create(ctx context.Context, booking *Booking) error
	GetByID(ctx context.Context, id string) (*Booking, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

// Service Interface (Use Cases)
type BookingService interface {
	CreateBooking(ctx context.Context, userID int64, eventID int64, quantity int32) (*Booking, string, error) // Returns Booking & PaymentURL
	GetBookingDetail(ctx context.Context, bookingID string, userID int64) (*Booking, error)
}

// Event Provider Interface
type EventProvider interface {
	CheckAvailability(ctx context.Context, eventID int64, quantity int32) (isAvailable bool, unitPrice float64, err error)
}
