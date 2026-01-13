package repository

import (
	"context"
	"errors"
	"time"

	"github.com/baobei23/e-ticket/services/payment-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var queryTimeoutDuration = 5 * time.Second

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) domain.PaymentRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, p *domain.Payment) error {
	query := `
		INSERT INTO payments (id, booking_id, user_id, amount, currency, status, stripe_id, payment_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := r.db.Exec(ctx, query,
		p.ID, p.BookingID, p.UserID, p.Amount, p.Currency,
		p.Status, p.StripeID, p.PaymentURL,
	)
	return err
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus) error {
	query := `UPDATE payments SET status = $1 WHERE id = $2`
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	tag, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("payment not found")
	}
	return nil
}

func (r *PostgresRepository) GetByBookingID(ctx context.Context, bookingID string) (*domain.Payment, error) {
	query := `
		SELECT id, booking_id, user_id, amount, currency, status, stripe_id, payment_url 
		FROM payments WHERE booking_id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var p domain.Payment
	err := r.db.QueryRow(ctx, query, bookingID).Scan(
		&p.ID, &p.BookingID, &p.UserID, &p.Amount, &p.Currency,
		&p.Status, &p.StripeID, &p.PaymentURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &p, nil
}
