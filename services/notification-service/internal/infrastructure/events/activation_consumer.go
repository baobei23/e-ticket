package events

import (
	"encoding/json"
	"log"

	"github.com/baobei23/e-ticket/services/notification-service/internal/domain"
	"github.com/baobei23/e-ticket/shared/contracts"
	"github.com/baobei23/e-ticket/shared/messaging"
)

type UserActivationConsumer struct {
	mq  *messaging.RabbitMQClient
	svc domain.NotificationService
}

func NewUserActivationConsumer(mq *messaging.RabbitMQClient, svc domain.NotificationService) *UserActivationConsumer {
	return &UserActivationConsumer{
		mq:  mq,
		svc: svc,
	}
}

func (c *UserActivationConsumer) Start() {
	if err := c.mq.Consume(contracts.QueueUserActivationRequested, c.handleMessage); err != nil {
		log.Printf("Failed to start consumer: %v", err)
	}
}

func (c *UserActivationConsumer) handleMessage(msg contracts.AmqpMessage) error {
	switch msg.EventName {
	case "UserActivationRequested":
		return c.processActivationRequested(msg.Payload)
	default:
		log.Printf("Unknown event: %s", msg.EventName)
		return nil
	}
}

func (c *UserActivationConsumer) processActivationRequested(payload []byte) error {
	var event contracts.UserActivationRequestedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	log.Printf("Send activation email to %s", event.Email)
	return c.svc.SendActivationEmail(event.Email, event.Token, event.ExpiresAt)
}
