package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	"github.com/baobei23/e-ticket/shared/contracts"
	"github.com/baobei23/e-ticket/shared/messaging"
)

type ExpiryConsumer struct {
	mq  *messaging.RabbitMQClient
	svc domain.BookingService
}

func NewExpiryConsumer(mq *messaging.RabbitMQClient, svc domain.BookingService) *ExpiryConsumer {
	return &ExpiryConsumer{
		mq:  mq,
		svc: svc,
	}
}

func (c *ExpiryConsumer) Start() {
	err := c.mq.Consume(contracts.QueueBookingExpiryProcess, c.handleMessage)
	if err != nil {
		log.Printf("Failed to start expiry consumer: %v", err)
	}
}

func (c *ExpiryConsumer) handleMessage(msg contracts.AmqpMessage) error {
	switch msg.EventName {
	case "BookingExpired":
		return c.processBookingExpiry(msg.Payload)
	default:
		return nil
	}
}

func (c *ExpiryConsumer) processBookingExpiry(payload []byte) error {
	var event contracts.BookingExpiryEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	log.Printf("Processing booking expiry for %s", event.BookingID)

	err := c.svc.ExpireBooking(context.Background(), event.BookingID)
	if err != nil {
		log.Printf("Failed to expire booking %s: %v", event.BookingID, err)
	}

	return nil
}
