package grpc

import (
	"context"
	"errors"

	"github.com/baobei23/e-ticket/services/auth-service/internal/domain"
	pb "github.com/baobei23/e-ticket/shared/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	svc domain.AuthService
}

func NewAuthHandler(server *grpc.Server, svc domain.AuthService) *AuthHandler {
	handler := &AuthHandler{svc: svc}
	pb.RegisterAuthServiceServer(server, handler)
	return handler
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, expiresIn, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotActive):
			return nil, status.Error(codes.PermissionDenied, "user not active")
		case errors.Is(err, domain.ErrInvalidCreds):
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userID, token, err := h.svc.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{UserId: userID, ActivationToken: token}, nil
}

func (h *AuthHandler) Activate(ctx context.Context, req *pb.ActivateRequest) (*pb.ActivateResponse, error) {
	if err := h.svc.Activate(ctx, req.Token); err != nil {
		return nil, err
	}
	return &pb.ActivateResponse{}, nil
}

func (h *AuthHandler) ResendActivation(ctx context.Context, req *pb.ResendActivationRequest) (*pb.ResendActivationResponse, error) {
	if err := h.svc.ResendActivation(ctx, req.Email); err != nil {
		return nil, err
	}
	return &pb.ResendActivationResponse{}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, err := h.svc.ValidateToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	// Optional: Fetch user detail untuk return email dll
	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
