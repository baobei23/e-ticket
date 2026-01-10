package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
)

type inMemBookingRepository struct {
	data map[string]*domain.Booking
	mu   sync.RWMutex
}

func NewInMemoryBookingRepository() domain.BookingRepository {
	return &inMemBookingRepository{
		data: make(map[string]*domain.Booking),
	}
}

func (r *inMemBookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[booking.ID] = booking
	return nil
}

func (r *inMemBookingRepository) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if booking, ok := r.data[id]; ok {
		return booking, nil
	}
	return nil, errors.New("booking not found")
}

func (r *inMemBookingRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if booking, ok := r.data[id]; ok {
		booking.Status = status
		return nil
	}
	return errors.New("booking not found")
}
