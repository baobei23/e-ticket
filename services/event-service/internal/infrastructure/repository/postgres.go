package repository

import (
	"context"
	"errors"
	"time"

	"github.com/baobei23/e-ticket/services/event-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

var queryTimeoutDuration = 5 * time.Second

func NewPostgresRepository(db *pgxpool.Pool) domain.EventRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAll(ctx context.Context, page, limit int) ([]*domain.Event, int64, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, name, description, location, start_time, end_time, total_seats, available_seats, price
		FROM events
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		var e domain.Event

		err := rows.Scan(
			&e.ID, &e.Name, &e.Description, &e.Location,
			&e.StartTime, &e.EndTime, &e.TotalSeats,
			&e.AvailableSeats, &e.Price,
		)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, &e)
	}

	var total int64
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM events").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*domain.Event, error) {
	query := `
		SELECT id, name, description, location, start_time, end_time, total_seats, available_seats, price
		FROM events
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var e domain.Event
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.Name, &e.Description, &e.Location,
		&e.StartTime, &e.EndTime, &e.TotalSeats,
		&e.AvailableSeats, &e.Price,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}

	return &e, nil
}

func (r *PostgresRepository) ReduceStock(ctx context.Context, eventID int64, quantity int32) error {
	// Atomic Update
	query := `
		UPDATE events
		SET available_seats = available_seats - $1
		WHERE id = $2 AND available_seats >= $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	tag, err := r.db.Exec(ctx, query, quantity, eventID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("insufficient seats or event not found")
	}

	return nil
}
