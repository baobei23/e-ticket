package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/baobei23/e-ticket/services/event-service/internal/domain"
	pb "github.com/baobei23/e-ticket/shared/proto/event"
)

type inMemRepository struct {
	data []*pb.Event
	mu   sync.RWMutex
}

func NewInMemoryRepository() domain.EventRepository {
	// Seed data dummy
	events := []*pb.Event{
		{
			Id:             1,
			Name:           "Konser Coldplay - Music of the Spheres",
			Description:    "World Tour 2025 di Jakarta",
			Location:       "Gelora Bung Karno",
			StartTime:      time.Now().Add(24 * time.Hour).Unix(),
			EndTime:        time.Now().Add(28 * time.Hour).Unix(),
			TotalSeats:     50000,
			AvailableSeats: 1000,
			Price:          3500000,
		},
		{
			Id:             2,
			Name:           "Tech Conference 2026",
			Description:    "Konferensi teknologi terbesar di Asia",
			Location:       "JCC Senayan",
			StartTime:      time.Now().Add(48 * time.Hour).Unix(),
			EndTime:        time.Now().Add(56 * time.Hour).Unix(),
			TotalSeats:     500,
			AvailableSeats: 500,
			Price:          150000,
		},
	}

	return &inMemRepository{
		data: events,
	}
}

func (r *inMemRepository) GetAll(ctx context.Context) ([]*pb.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data, nil
}

func (r *inMemRepository) GetByID(ctx context.Context, id int64) (*pb.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, e := range r.data {
		if e.Id == id {
			return e, nil
		}
	}
	return nil, errors.New("event not found")
}

func (r *inMemRepository) ReduceStock(ctx context.Context, eventID int64, quantity int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, e := range r.data {
		if e.Id == eventID {
			if e.AvailableSeats < quantity {
				return errors.New("insufficient seats")
			}
			e.AvailableSeats -= quantity
			return nil
		}
	}
	return errors.New("event not found")
}
