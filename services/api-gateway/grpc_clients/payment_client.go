package grpc_clients

import (
	"github.com/baobei23/e-ticket/shared/env"
	paymentpb "github.com/baobei23/e-ticket/shared/proto/payment"
	"github.com/baobei23/e-ticket/shared/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentServiceClient struct {
	Client paymentpb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentServiceClient() (*PaymentServiceClient, error) {
	addr := env.GetString("PAYMENT_SERVICE_ADDRESS", "payment-service:50053")

	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(addr, dialOptions...)
	if err != nil {
		return nil, err
	}

	client := paymentpb.NewPaymentServiceClient(conn)

	return &PaymentServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *PaymentServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
