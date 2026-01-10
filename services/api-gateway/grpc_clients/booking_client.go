package grpc_clients

import (
	"os"

	"github.com/baobei23/e-ticket/shared/proto/booking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookingServiceClient struct {
	Client booking.BookingServiceClient
	conn   *grpc.ClientConn
}

func NewBookingServiceClient() (*BookingServiceClient, error) {
	svcAddr := os.Getenv("BOOKING_SERVICE_ADDRESS")
	if svcAddr == "" {
		svcAddr = "booking-service:50052"
	}

	conn, err := grpc.NewClient(svcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := booking.NewBookingServiceClient(conn)

	return &BookingServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *BookingServiceClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
