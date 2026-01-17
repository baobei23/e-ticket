package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/baobei23/e-ticket/services/payment-service/internal/domain"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type PaymentService struct {
	repo      domain.PaymentRepository
	gateway   domain.PaymentGateway
	publisher domain.PaymentPublisher
}

func NewPaymentService(
	repo domain.PaymentRepository,
	gateway domain.PaymentGateway,
	publisher domain.PaymentPublisher) *PaymentService {
	return &PaymentService{
		repo:      repo,
		gateway:   gateway,
		publisher: publisher,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, bookingID string, userID int64, amount float64, unitPrice float64, quantity int32) (*domain.Payment, error) {

	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}
	if unitPrice <= 0 || quantity <= 0 {
		return nil, errors.New("invalid unit price or quantity")
	}

	paymentID := uuid.New().String()
	payment := &domain.Payment{
		ID:        paymentID,
		BookingID: bookingID,
		UserID:    userID,
		Amount:    amount,
		UnitPrice: unitPrice,
		Quantity:  quantity,
		Currency:  "idr",
		Status:    domain.PaymentStatusPending,
	}

	stripeID, paymentURL, err := s.gateway.CreateSession(ctx, payment)
	if err != nil {
		return nil, err
	}

	payment.StripeID = stripeID
	payment.PaymentURL = paymentURL

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) HandleWebhook(ctx context.Context, payload []byte, sigHeader string) error {
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		return fmt.Errorf("STRIPE_WEBHOOK_SECRET is not set")
	}

	// Verify signature from Stripe
	event, err := webhook.ConstructEventWithOptions(payload, sigHeader, webhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	log.Printf("Received Webhook Event: %s", event.Type)

	// Handle event type
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return fmt.Errorf("failed to unmarshal session: %w", err)
		}

		// Get payment_id from metadata that we set when creating session
		paymentID := session.Metadata["payment_id"]
		if paymentID == "" {
			return fmt.Errorf("payment_id not found in metadata")
		}

		log.Printf("Payment success for Payment ID: %s", paymentID)

		// Update status payment
		if err := s.repo.UpdateStatus(ctx, paymentID, domain.PaymentStatusSuccess); err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		paymentData, err := s.repo.GetByBookingID(ctx, session.Metadata["booking_id"])
		if err != nil {
			log.Printf("Failed to get payment data: %v", err)
			return fmt.Errorf("failed to get payment data: %w", err)
		}

		if err := s.publisher.PublishPaymentSuccess(ctx, paymentData); err != nil {
			log.Printf("Failed to publish PaymentSuccess: %v", err)
		}
	case "checkout.session.expired":
		// Handle expired...
	}

	return nil
}
