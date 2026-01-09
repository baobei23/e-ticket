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

	page := int(req.Pagination.GetPage())
	if page == 0 {
		page = 1
	}
	limit := int(req.Pagination.GetLimit())
	if limit == 0 {
		limit = 10
	}

	events, totalItems, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int32(0)
	if limit > 0 {
		totalPages = int32((totalItems + int64(limit) - 1) / int64(limit))
	}

	return &pb.GetEventsResponse{
		Events: events,
		Meta: &event.PaginationMetadata{
			CurrentPage: int32(page),
			PageLimit:   int32(limit),
			TotalItems:  int32(totalItems),
			TotalPages:  totalPages,
		},
	}, nil
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

func (s *EventService) GetEventDetail(ctx context.Context, req *pb.GetEventDetailRequest) (*pb.GetEventDetailResponse, error) {
	event, err := s.repo.GetByID(ctx, req.EventId)
	if err != nil {
		return nil, err
	}

	return &pb.GetEventDetailResponse{
		Event: event,
	}, nil
}
