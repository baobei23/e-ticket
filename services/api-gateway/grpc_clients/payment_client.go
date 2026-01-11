package grpc_clients

import (
	"os"

	"github.com/baobei23/e-ticket/shared/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentServiceClient struct {
	Client payment.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentServiceClient() (*PaymentServiceClient, error) {
	svcAddr := os.Getenv("PAYMENT_SERVICE_ADDRESS")
	if svcAddr == "" {
		svcAddr = "payment-service:50053"
	}

	conn, err := grpc.NewClient(svcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentServiceClient(conn)

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
