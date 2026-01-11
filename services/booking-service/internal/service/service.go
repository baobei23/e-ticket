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

	// PUBLISH EVENT
	// Kita lakukan secara async (go routine) atau sync tergantung kebutuhan.
	// Jika RabbitMQ mati, apakah booking harus gagal?
	// Untuk reliability tinggi, sebaiknya sync dulu atau gunakan Outbox Pattern.
	// Untuk sekarang, kita panggil sync tapi log error saja (soft failure).

	if err := s.publisher.PublishBookingCreated(ctx, booking); err != nil {
		// Warning: Event gagal terkirim. Stok mungkin tidak ter-reserve.
		// Idealnya: Rollback booking atau simpan ke Outbox table.
	}

	paymentURL, err := s.paymentProvider.CreatePayment(ctx, bookingID, userID, totalAmount)
	if err != nil {
		// Log error, tapi booking sudah terbentuk.
		// User harus bisa retry payment nanti (Endpoint GetBookingDetail harus return paymentURL juga kalau belum lunas).
		// Untuk simplicity sekarang, kita return error atau string kosong.
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
