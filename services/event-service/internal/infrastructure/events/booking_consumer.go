package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/baobei23/e-ticket/services/event-service/internal/domain"
	"github.com/baobei23/e-ticket/shared/contracts"
	"github.com/baobei23/e-ticket/shared/messaging"
)

type EventConsumer struct {
	mq  *messaging.RabbitMQClient
	svc domain.EventService
}

func NewEventConsumer(mq *messaging.RabbitMQClient, svc domain.EventService) *EventConsumer {
	return &EventConsumer{
		mq:  mq,
		svc: svc,
	}
}

func (c *EventConsumer) Start() {

	err := c.mq.Consume(contracts.QueueBookingCreated, c.handleMessage)
	if err != nil {
		log.Printf("Failed to start consumer: %v", err)
	}
}

func (c *EventConsumer) handleMessage(msg contracts.AmqpMessage) error {
	log.Printf("Received event: %s", msg.EventName)

	switch msg.EventName {
	case "BookingCreated":
		return c.processBookingCreated(msg.Payload)
	default:
		log.Printf("Unknown event: %s", msg.EventName)
		return nil
	}
}

func (c *EventConsumer) processBookingCreated(payload []byte) error {
	var event contracts.BookingCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	log.Printf("Processing booking %s for event %d (Qty: %d)", event.BookingID, event.EventID, event.Quantity)

	err := c.svc.ReduceStock(context.Background(), event.EventID, event.Quantity)
	if err != nil {
		log.Printf("Error reducing stock: %v", err)
		return err
	}

	log.Printf("Successfully reserved stock for booking %s", event.BookingID)

	return nil
}
