package domain

import (
	"context"
	"errors"
	"time"
)

var (
	QueryTimeoutDuration = 5 * time.Second
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotActive     = errors.New("user not active")
	ErrInvalidCreds      = errors.New("invalid credentials")
	ErrEmailExists       = errors.New("email already exists")
	ErrInvalidToken      = errors.New("invalid or expired token")
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  []byte    `json:"-"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User, token string, expiry time.Duration) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	ActivateByToken(ctx context.Context, token string) error
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (int64, string, error)
	Login(ctx context.Context, email, password string) (string, int64, error)
	ValidateToken(ctx context.Context, token string) (int64, error)
	Activate(ctx context.Context, token string) error
}
