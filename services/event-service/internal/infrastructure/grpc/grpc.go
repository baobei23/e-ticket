package grpc

import (
	"context"

	"github.com/baobei23/e-ticket/services/event-service/internal/domain"
	"github.com/baobei23/e-ticket/shared/proto/event"
	pb "github.com/baobei23/e-ticket/shared/proto/event"
	"google.golang.org/grpc"
)

type EventHandler struct {
	pb.UnimplementedEventServiceServer
	service domain.EventService
}

func NewEventHandler(server *grpc.Server, service domain.EventService) *EventHandler {
	handler := &EventHandler{
		service: service,
	}
	pb.RegisterEventServiceServer(server, handler)
	return handler
}

func (h *EventHandler) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	page := int(req.Pagination.GetPage())
	if page == 0 {
		page = 1
	}
	limit := int(req.Pagination.GetLimit())
	if limit == 0 {
		limit = 10
	}

	events, totalItems, err := h.service.GetEvents(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	var pbEvents []*pb.Event
	for _, e := range events {
		pbEvents = append(pbEvents, e.ToProto())
	}

	totalPages := int32(0)
	if limit > 0 {
		totalPages = int32((totalItems + int64(limit) - 1) / int64(limit))
	}

	return &pb.GetEventsResponse{
		Events: pbEvents,
		Meta: &event.PaginationMetadata{
			CurrentPage: int32(page),
			PageLimit:   int32(limit),
			TotalItems:  int32(totalItems),
			TotalPages:  totalPages,
		},
	}, nil
}

func (h *EventHandler) GetEventDetail(ctx context.Context, req *pb.GetEventDetailRequest) (*pb.GetEventDetailResponse, error) {
	event, err := h.service.GetEventDetail(ctx, req.EventId)
	if err != nil {
		return nil, err
	}

	return &pb.GetEventDetailResponse{
		Event: event.ToProto(),
	}, nil
}

func (h *EventHandler) CheckAvailability(ctx context.Context, req *pb.CheckAvailabilityRequest) (*pb.CheckAvailabilityResponse, error) {
	available, price, err := h.service.CheckAvailability(ctx, req.EventId, req.Quantity)
	if err != nil {
		return &pb.CheckAvailabilityResponse{IsAvailable: false}, nil
	}

	return &pb.CheckAvailabilityResponse{
		IsAvailable: available,
		UnitPrice:   price,
	}, nil
}
