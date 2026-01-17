package grpc

import (
	"context"

	"github.com/baobei23/e-ticket/services/payment-service/internal/domain"
	paymentpb "github.com/baobei23/e-ticket/shared/proto/payment"
	"google.golang.org/grpc"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	svc domain.PaymentService
}

func NewPaymentHandler(server *grpc.Server, svc domain.PaymentService) *PaymentHandler {
	handler := &PaymentHandler{svc: svc}
	paymentpb.RegisterPaymentServiceServer(server, handler)
	return handler
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	payment, err := h.svc.CreatePayment(ctx, req.BookingId, req.UserId, req.Amount, req.UnitPrice, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &paymentpb.CreatePaymentResponse{
		PaymentId:  payment.ID,
		PaymentUrl: payment.PaymentURL,
		Status:     string(payment.Status),
	}, nil
}

func (h *PaymentHandler) HandleWebhook(ctx context.Context, req *paymentpb.HandleWebhookRequest) (*paymentpb.HandleWebhookResponse, error) {
	err := h.svc.HandleWebhook(ctx, req.Payload, req.Signature)
	if err != nil {
		return nil, err
	}

	return &paymentpb.HandleWebhookResponse{Success: true}, nil
}
