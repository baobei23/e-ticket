package service

import (
	"context"
	"errors"
	"time"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	"github.com/google/uuid"
)

type BookingService struct {
	repo          domain.BookingRepository
	eventProvider domain.EventProvider
}

func NewBookingService(repo domain.BookingRepository, eventProvider domain.EventProvider) domain.BookingService {
	return &BookingService{
		repo:          repo,
		eventProvider: eventProvider,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID int64, eventID int64, quantity int32) (*domain.Booking, string, error) {
	// 1. Panggil Interface (Bersih dari protobuf struct)
	isAvailable, unitPrice, err := s.eventProvider.CheckAvailability(ctx, eventID, quantity)
	if err != nil {
		return nil, "", err
	}

	if !isAvailable {
		return nil, "", errors.New("insufficient seats or event not found")
	}

	// 2. Logic selanjutnya sama...
	totalAmount := unitPrice * float64(quantity)
	bookingID := uuid.New().String()

	booking := &domain.Booking{
		ID:          bookingID,
		UserID:      userID,
		EventID:     eventID,
		Quantity:    quantity,
		TotalAmount: totalAmount,
		Status:      domain.StatusPending,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, booking); err != nil {
		return nil, "", err
	}

	paymentURL := "https://payment-gateway.com/pay/" + bookingID
	return booking, paymentURL, nil
}
func (s *BookingService) GetBookingDetail(ctx context.Context, bookingID string, userID int64) (*domain.Booking, error) {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.UserID != userID {
		return nil, errors.New("unauthorized access to booking")
	}

	return booking, nil
}
