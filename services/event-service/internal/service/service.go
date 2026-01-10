package service

import (
	"context"

	"github.com/baobei23/e-ticket/services/event-service/internal/domain"
)

type EventService struct {
	repo domain.EventRepository
}

func NewEventService(repo domain.EventRepository) domain.EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) GetEvents(ctx context.Context, page, limit int) ([]*domain.Event, int64, error) {
	return s.repo.GetAll(ctx, page, limit)
}

func (s *EventService) GetEventDetail(ctx context.Context, id int64) (*domain.Event, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EventService) CheckAvailability(ctx context.Context, eventID int64, quantity int32) (bool, float64, error) {
	event, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return false, 0, err
	}

	isAvailable := event.AvailableSeats >= quantity

	return isAvailable, event.Price, nil
}
