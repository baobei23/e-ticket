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

func NewBookingService(repo domain.BookingRepository, eventProvider domain.EventProvider, publisher domain.BookingPublisher, paymentProvider domain.PaymentProvider) domain.BookingService {
	return &BookingService{
		repo:            repo,
		eventProvider:   eventProvider,
		publisher:       publisher,
		paymentProvider: paymentProvider,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID int64, eventID int64, quantity int32) (*domain.Booking, string, error) {

	isAvailable, unitPrice, err := s.eventProvider.CheckAvailability(ctx, eventID, quantity)
	if err != nil {
		return nil, "", err
	}

	if !isAvailable {
		return nil, "", errors.New("insufficient seats")
	}

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

	// TODO: outbox pattern

	if err := s.publisher.PublishBookingCreated(ctx, booking); err != nil {
		// TODO: outbox pattern
	}

	paymentURL, err := s.paymentProvider.CreatePayment(ctx, bookingID, userID, totalAmount, unitPrice, quantity)
	if err != nil {
		// Log error, but booking is created
		// User should be able to retry payment later (Endpoint GetBookingDetail should return paymentURL also if not paid),
		// for simplicity, we return error or empty string.
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
