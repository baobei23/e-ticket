package repository

import (
	"context"
	"errors"
	"time"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) domain.BookingRepository {
	return &PostgresRepository{db: db}
}

var queryTimeoutDuration = 5 * time.Second

func (r *PostgresRepository) Create(ctx context.Context, booking *domain.Booking) error {
	query := `
		INSERT INTO bookings (id, user_id, event_id, quantity, total_amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	_, err := r.db.Exec(ctx, query,
		booking.ID, booking.UserID, booking.EventID,
		booking.Quantity, booking.TotalAmount, booking.Status, booking.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	query := `SELECT id, user_id, event_id, quantity, total_amount, status, created_at FROM bookings WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	var b domain.Booking
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.UserID, &b.EventID, &b.Quantity,
		&b.TotalAmount, &b.Status, &b.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}
	return &b, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE bookings SET status = $1 WHERE id = $2`
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	tag, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("booking not found")
	}
	return nil
}
