package domain

import (
	"context"

	pb "github.com/baobei23/e-ticket/shared/proto/event"
)

type EventRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]*pb.Event, int64, error)
	GetByID(ctx context.Context, id int64) (*pb.Event, error)
	ReduceStock(ctx context.Context, eventID int64, quantity int32) error
}

type EventService interface {
	GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error)
	GetEventDetail(ctx context.Context, req *pb.GetEventDetailRequest) (*pb.GetEventDetailResponse, error)
	CheckAvailability(ctx context.Context, req *pb.CheckAvailabilityRequest) (*pb.CheckAvailabilityResponse, error)
}
