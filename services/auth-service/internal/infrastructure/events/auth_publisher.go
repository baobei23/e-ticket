package events

import (
	"context"
	"time"

	"github.com/baobei23/e-ticket/shared/contracts"
	"github.com/baobei23/e-ticket/shared/messaging"
)

type UserActivationPublisher struct {
	mq *messaging.RabbitMQClient
}

func NewUserActivationPublisher(mq *messaging.RabbitMQClient) *UserActivationPublisher {
	return &UserActivationPublisher{mq: mq}
}

func (p *UserActivationPublisher) Publish(ctx context.Context, userID int64, email, token string, expiresAt time.Time) error {
	payload := contracts.UserActivationRequestedEvent{
		UserID:    userID,
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return p.mq.Publish(ctx, "UserActivationRequested", payload)
}
