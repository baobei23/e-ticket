package service

import (
	"context"

	"github.com/baobei23/e-ticket/services/event-service/internal/domain"

	"github.com/baobei23/e-ticket/shared/proto/event"
	pb "github.com/baobei23/e-ticket/shared/proto/event"
)

type EventService struct {
	pb.UnimplementedEventServiceServer
	repo domain.EventRepository
}

func NewEventService(repo domain.EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetEventsResponse{
		Events: events,
		Meta: &event.PaginationMetadata{
			TotalItems: int32(len(events)),
			// Simplifikasi: abaikan pagination logic dulu
		},
	}, nil
}

func (s *EventService) GetEventDetail(ctx context.Context, req *pb.GetEventDetailRequest) (*pb.GetEventDetailResponse, error) {
	event, err := s.repo.GetByID(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return &pb.GetEventDetailResponse{Event: event}, nil
}

func (s *EventService) CheckAvailability(ctx context.Context, req *pb.CheckAvailabilityRequest) (*pb.CheckAvailabilityResponse, error) {
	event, err := s.repo.GetByID(ctx, req.EventId)
	if err != nil {
		return &pb.CheckAvailabilityResponse{IsAvailable: false}, nil
	}

	isAvailable := event.AvailableSeats >= req.Quantity

	return &pb.CheckAvailabilityResponse{
		IsAvailable: isAvailable,
		UnitPrice:   event.Price,
	}, nil
}
