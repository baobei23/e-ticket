package gateway

import (
	"context"
	"fmt"

	"github.com/baobei23/e-ticket/services/payment-service/internal/domain"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type StripeGateway struct {
	apiKey     string
	successURL string
	cancelURL  string
}

func NewStripeGateway(apiKey, successURL, cancelURL string) *StripeGateway {
	stripe.Key = apiKey
	return &StripeGateway{
		apiKey:     apiKey,
		successURL: successURL,
		cancelURL:  cancelURL,
	}
}

func (s *StripeGateway) CreateSession(ctx context.Context, payment *domain.Payment) (string, string, error) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card", "grabpay"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(payment.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf("Booking #%s", payment.BookingID)),
					},
					UnitAmount: stripe.Int64(int64(payment.Amount * 100)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.successURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.cancelURL + "?session_id={CHECKOUT_SESSION_ID}"),

		Metadata: map[string]string{
			"booking_id": payment.BookingID,
			"user_id":    fmt.Sprintf("%d", payment.UserID),
			"payment_id": payment.ID,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return "", "", fmt.Errorf("stripe session creation failed: %w", err)
	}

	return sess.ID, sess.URL, nil
}
