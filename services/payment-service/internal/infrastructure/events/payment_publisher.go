package events

import (
	"context"

	"github.com/baobei23/e-ticket/services/payment-service/internal/domain"
	"github.com/baobei23/e-ticket/shared/contracts"
	"github.com/baobei23/e-ticket/shared/messaging"
)

type PaymentEventPublisher struct {
	mq *messaging.RabbitMQClient
}

func NewPaymentEventPublisher(mq *messaging.RabbitMQClient) domain.PaymentPublisher {
	return &PaymentEventPublisher{mq: mq}
}

func (p *PaymentEventPublisher) PublishPaymentSuccess(ctx context.Context, payment *domain.Payment) error {
	payload := contracts.PaymentSuccessEvent{
		PaymentID: payment.ID,
		BookingID: payment.BookingID,
		Amount:    payment.Amount,
	}

	return p.mq.Publish(ctx, "PaymentSuccess", payload)
}

func (p *PaymentEventPublisher) PublishPaymentFailed(ctx context.Context, payment *domain.Payment, reason string) error {
	payload := contracts.PaymentFailedEvent{
		BookingID: payment.BookingID,
		Reason:    reason,
	}

	return p.mq.Publish(ctx, "PaymentFailed", payload)
}
