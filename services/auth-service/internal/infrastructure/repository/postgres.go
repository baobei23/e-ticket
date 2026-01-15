package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/baobei23/e-ticket/services/auth-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) domain.UserRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, user *domain.User, token string, expiry time.Duration) error {
	return withTx(ctx, r.db, func(tx pgx.Tx) error {
		insertUser := `INSERT INTO users (email, password, is_active, created_at) VALUES ($1, $2, $3, $4) RETURNING id`

		ctx, cancel := context.WithTimeout(ctx, domain.QueryTimeoutDuration)
		defer cancel()

		if err := tx.QueryRow(ctx, insertUser, user.Email, user.Password, false, user.CreatedAt).Scan(&user.ID); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return domain.ErrEmailExists
			}
			return err
		}

		insertToken := `INSERT INTO user_activation_tokens (token, user_id, expiry) VALUES ($1, $2, $3)`
		_, err := tx.Exec(ctx, insertToken, hashToken(token), user.ID, time.Now().Add(expiry))
		return err
	})
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password, is_active, created_at FROM users WHERE email = $1 AND is_active = true`

	ctx, cancel := context.WithTimeout(ctx, domain.QueryTimeoutDuration)
	defer cancel()

	var u domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password, &u.IsActive, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *PostgresRepository) ActivateByToken(ctx context.Context, token string) error {
	return withTx(ctx, r.db, func(tx pgx.Tx) error {
		selectUser := `SELECT user_id FROM user_activation_tokens WHERE token = $1 AND expiry > NOW()`
		ctx, cancel := context.WithTimeout(ctx, domain.QueryTimeoutDuration)
		defer cancel()

		var userID int64
		if err := tx.QueryRow(ctx, selectUser, hashToken(token)).Scan(&userID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrInvalidToken
			}
			return err
		}

		updateUser := `UPDATE users SET is_active = true WHERE id = $1`
		if _, err := tx.Exec(ctx, updateUser, userID); err != nil {
			return err
		}

		deleteTokens := `DELETE FROM user_activation_tokens WHERE user_id = $1`
		_, err := tx.Exec(ctx, deleteTokens, userID)
		return err
	})
}

// helper functions
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func withTx(ctx context.Context, db *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
