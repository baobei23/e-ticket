package clients

import (
	"context"
	"fmt"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	"github.com/baobei23/e-ticket/shared/env"
	paymentpb "github.com/baobei23/e-ticket/shared/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentGRPCClient struct {
	client paymentpb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentGRPCClient() (domain.PaymentProvider, error) {
	addr := env.GetString("PAYMENT_SERVICE_ADDRESS", "payment-service:50053")
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}

	client := paymentpb.NewPaymentServiceClient(conn)

	return &PaymentGRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *PaymentGRPCClient) CreatePayment(ctx context.Context, bookingID string, userID int64, amount float64) (string, error) {
	resp, err := c.client.CreatePayment(ctx, &paymentpb.CreatePaymentRequest{
		BookingId: bookingID,
		UserId:    userID,
		Amount:    amount,
		Currency:  "idr",
	})
	if err != nil {
		return "", err
	}
	return resp.PaymentUrl, nil
}
