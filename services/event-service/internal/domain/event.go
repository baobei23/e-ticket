package domain

import (
	"context"
	"time"

	pb "github.com/baobei23/e-ticket/shared/proto/event"
)

type Event struct {
	ID             int64
	Name           string
	Description    string
	Location       string
	StartTime      time.Time
	EndTime        time.Time
	TotalSeats     int32
	AvailableSeats int32
	Price          float64
}

func (e *Event) ToProto() *pb.Event {
	return &pb.Event{
		Id:             e.ID,
		Name:           e.Name,
		Description:    e.Description,
		Location:       e.Location,
		StartTime:      e.StartTime.Unix(),
		EndTime:        e.EndTime.Unix(),
		TotalSeats:     e.TotalSeats,
		AvailableSeats: e.AvailableSeats,
		Price:          e.Price,
	}
}

type EventRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]*Event, int64, error)
	GetByID(ctx context.Context, id int64) (*Event, error)
	ReduceStock(ctx context.Context, eventID int64, quantity int32) error
}

type EventService interface {
	GetEvents(ctx context.Context, page, limit int) ([]*Event, int64, error)
	GetEventDetail(ctx context.Context, id int64) (*Event, error)
	CheckAvailability(ctx context.Context, eventID int64, quantity int32) (bool, float64, error)
	ReduceStock(ctx context.Context, eventID int64, quantity int32) error
}
