package service

import (
	"context"
	"os"
	"time"

	"github.com/baobei23/e-ticket/services/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo        domain.UserRepository
	jwtSecret   []byte
	tokenExpiry time.Duration
	publisher   domain.UserActivationPublisher
}

const activationTokenTTL = 30 * time.Minute

func NewAuthService(repo domain.UserRepository, publisher domain.UserActivationPublisher) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey" // Default dev
	}
	return &AuthService{
		repo:        repo,
		jwtSecret:   []byte(secret),
		tokenExpiry: 24 * time.Hour,
		publisher:   publisher,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (int64, string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, "", err
	}

	user := &domain.User{
		Email:     email,
		Password:  hashedPass,
		CreatedAt: time.Now(),
	}

	plainToken := uuid.New().String()
	expiry, err := s.repo.Create(ctx, user, plainToken, activationTokenTTL)
	if err != nil {
		return 0, "", err
	}

	if err := s.publisher.Publish(ctx, user.ID, user.Email, plainToken, expiry); err != nil {
		return 0, "", err
	}

	return user.ID, plainToken, nil
}

func (s *AuthService) Activate(ctx context.Context, token string) error {
	return s.repo.ActivateByToken(ctx, token)
}

func (s *AuthService) ResendActivation(ctx context.Context, email string) error {
	user, err := s.repo.GetByEmailAnyStatus(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil
		}
		return err
	}

	if user.IsActive {
		return nil
	}

	plainToken := uuid.New().String()
	expiry, err := s.repo.UpsertActivationToken(ctx, user.ID, plainToken, activationTokenTTL)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, user.ID, user.Email, plainToken, expiry)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, int64, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", 0, domain.ErrUserNotFound
	}

	// Compare Password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return "", 0, domain.ErrInvalidCreds
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(s.tokenExpiry.Seconds()), nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, domain.ErrInvalidToken
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, domain.ErrInvalidToken
	}

	return int64(userIDFloat), nil
}
