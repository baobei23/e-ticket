package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/baobei23/e-ticket/services/event-service/internal/domain"
)

type inMemRepository struct {
	data []*domain.Event
	mu   sync.RWMutex
}

func NewInMemoryRepository() domain.EventRepository {
	// Seed data dummy
	events := []*domain.Event{
		{
			ID:             1,
			Name:           "Konser Coldplay - Music of the Spheres",
			Description:    "World Tour 2025 di Jakarta",
			Location:       "Gelora Bung Karno",
			StartTime:      time.Now().Add(24 * time.Hour),
			EndTime:        time.Now().Add(28 * time.Hour),
			TotalSeats:     50000,
			AvailableSeats: 1000,
			Price:          3500000,
		},
		{
			ID:             2,
			Name:           "Tech Conference 2026",
			Description:    "Konferensi teknologi terbesar di Asia",
			Location:       "JCC Senayan",
			StartTime:      time.Now().Add(48 * time.Hour),
			EndTime:        time.Now().Add(56 * time.Hour),
			TotalSeats:     500,
			AvailableSeats: 500,
			Price:          150000,
		},
	}

	return &inMemRepository{
		data: events,
	}
}

func (r *inMemRepository) GetAll(ctx context.Context, page, limit int) ([]*domain.Event, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	totalItems := int64(len(r.data))

	// Hitung Start & End Index
	start := (page - 1) * limit
	end := start + limit

	// Validasi Bounds
	if start >= len(r.data) {
		return []*domain.Event{}, totalItems, nil // Halaman kosong (out of range)
	}
	if end > len(r.data) {
		end = len(r.data)
	}

	// Return slice sesuai halaman
	return r.data[start:end], totalItems, nil
}

func (r *inMemRepository) GetByID(ctx context.Context, id int64) (*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, e := range r.data {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, errors.New("event not found")
}

func (r *inMemRepository) ReduceStock(ctx context.Context, eventID int64, quantity int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, e := range r.data {
		if e.ID == eventID {
			if e.AvailableSeats < quantity {
				return errors.New("insufficient seats")
			}
			e.AvailableSeats -= quantity
			return nil
		}
	}
	return errors.New("event not found")
}
