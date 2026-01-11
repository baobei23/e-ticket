package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/baobei23/e-ticket/services/payment-service/internal/domain"
)

type inMemPaymentRepository struct {
	data map[string]*domain.Payment
	mu   sync.RWMutex
}

func NewInMemoryPaymentRepository() domain.PaymentRepository {
	return &inMemPaymentRepository{
		data: make(map[string]*domain.Payment),
	}
}

func (r *inMemPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[payment.ID] = payment
	return nil
}

func (r *inMemPaymentRepository) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if payment, ok := r.data[id]; ok {
		payment.Status = status
		return nil
	}
	return errors.New("payment not found")
}

func (r *inMemPaymentRepository) GetByBookingID(ctx context.Context, bookingID string) (*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.data {
		if p.BookingID == bookingID {
			return p, nil
		}
	}
	return nil, errors.New("payment not found")
}
