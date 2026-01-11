package grpc

import (
	"context"

	"github.com/baobei23/e-ticket/services/payment-service/internal/domain"
	pb "github.com/baobei23/e-ticket/shared/proto/payment"
	"google.golang.org/grpc"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	svc domain.PaymentService
}

func NewPaymentHandler(server *grpc.Server, svc domain.PaymentService) *PaymentHandler {
	handler := &PaymentHandler{svc: svc}
	pb.RegisterPaymentServiceServer(server, handler)
	return handler
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	payment, err := h.svc.CreatePayment(ctx, req.BookingId, req.UserId, req.Amount)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		PaymentId:  payment.ID,
		PaymentUrl: payment.PaymentURL,
		Status:     string(payment.Status),
	}, nil
}

func (h *PaymentHandler) HandleWebhook(ctx context.Context, req *pb.HandleWebhookRequest) (*pb.HandleWebhookResponse, error) {
	err := h.svc.HandleWebhook(ctx, req.Payload, req.Signature)
	if err != nil {
		return nil, err
	}

	return &pb.HandleWebhookResponse{Success: true}, nil
}
