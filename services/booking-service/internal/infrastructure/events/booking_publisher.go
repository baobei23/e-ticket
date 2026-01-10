package events

import (
	"context"
	"log"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	"github.com/baobei23/e-ticket/shared/contracts"
	"github.com/baobei23/e-ticket/shared/messaging"
)

type BookingEventPublisher struct {
	mq *messaging.RabbitMQClient
}

func NewBookingEventPublisher(mq *messaging.RabbitMQClient) domain.BookingPublisher {
	return &BookingEventPublisher{
		mq: mq,
	}
}

func (p *BookingEventPublisher) PublishBookingCreated(ctx context.Context, booking *domain.Booking) error {
	payload := contracts.BookingCreatedEvent{
		BookingID: booking.ID,
		UserID:    booking.UserID,
		EventID:   booking.EventID,
		Quantity:  booking.Quantity,
		CreatedAt: booking.CreatedAt,
	}

	err := p.mq.Publish(ctx, "BookingCreated", payload)

	if err != nil {
		log.Printf("Failed to publish BookingCreated event: %v", err)
		return err
	}

	return nil
}
