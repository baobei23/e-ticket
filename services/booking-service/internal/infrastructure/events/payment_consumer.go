package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	"github.com/baobei23/e-ticket/shared/contracts"
	"github.com/baobei23/e-ticket/shared/messaging"
)

type PaymentConsumer struct {
	mq  *messaging.RabbitMQClient
	svc domain.BookingService
}

func NewPaymentConsumer(mq *messaging.RabbitMQClient, svc domain.BookingService) *PaymentConsumer {
	return &PaymentConsumer{
		mq:  mq,
		svc: svc,
	}
}

func (c *PaymentConsumer) Start() {
	err := c.mq.Consume(contracts.QueuePaymentSuccess, c.handleMessage)
	if err != nil {
		log.Printf("Failed to start payment consumer: %v", err)
	}
}

func (c *PaymentConsumer) handleMessage(msg contracts.AmqpMessage) error {
	switch msg.EventName {
	case "PaymentSuccess":
		return c.processPaymentSuccess(msg.Payload)
	default:
		return nil
	}
}

func (c *PaymentConsumer) processPaymentSuccess(payload []byte) error {
	var event contracts.PaymentSuccessEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	log.Printf("Payment success for booking %s. Confirming...", event.BookingID)

	return c.svc.ConfirmBooking(context.Background(), event.BookingID)
}
