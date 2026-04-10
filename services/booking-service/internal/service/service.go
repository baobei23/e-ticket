package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	"github.com/google/uuid"
)

type BookingService struct {
	repo            domain.BookingRepository
	eventProvider   domain.EventProvider
	publisher       domain.BookingPublisher
	paymentProvider domain.PaymentProvider
}

func NewBookingService(repo domain.BookingRepository, eventProvider domain.EventProvider, publisher domain.BookingPublisher, paymentProvider domain.PaymentProvider) *BookingService {
	return &BookingService{
		repo:            repo,
		eventProvider:   eventProvider,
		publisher:       publisher,
		paymentProvider: paymentProvider,
	}
}

// main flow of creating a booking
func (s *BookingService) CreateBooking(ctx context.Context, userID int64, eventID int64, quantity int32) (*domain.Booking, string, error) {

	bookingID := uuid.New().String()

	if err := s.eventProvider.ReserveSeat(ctx, eventID, quantity, bookingID); err != nil {
		return nil, "", fmt.Errorf("failed to reserve seat: %w", err)
	}

	_, unitPrice, err := s.eventProvider.CheckAvailability(ctx, eventID, quantity)
	if err != nil {
		s.eventProvider.ReleaseSeat(ctx, eventID, quantity, bookingID)
		return nil, "", fmt.Errorf("failed to check availability: %w", err)
	}

	totalAmount := unitPrice * float64(quantity)

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
		s.eventProvider.ReleaseSeat(ctx, eventID, quantity, bookingID)
		return nil, "", fmt.Errorf("failed to create booking: %w", err)
	}

	paymentURL, err := s.paymentProvider.CreatePayment(ctx, bookingID, userID, totalAmount, unitPrice, quantity)
	if err != nil {
		s.eventProvider.ReleaseSeat(ctx, eventID, quantity, bookingID)
		return nil, "", fmt.Errorf("failed to create payment session: %w", err)
	}

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

func (s *BookingService) ConfirmBooking(ctx context.Context, bookingID string) error {
	return s.repo.UpdateStatus(ctx, bookingID, domain.StatusConfirmed)
}

func (s *BookingService) FailBooking(ctx context.Context, bookingID string) error {
	return s.repo.UpdateStatus(ctx, bookingID, domain.StatusFailed)
}
