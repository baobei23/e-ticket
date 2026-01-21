package grpc_clients

import (
	"github.com/baobei23/e-ticket/shared/env"
	bookingpb "github.com/baobei23/e-ticket/shared/proto/booking"
	"github.com/baobei23/e-ticket/shared/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookingServiceClient struct {
	Client bookingpb.BookingServiceClient
	conn   *grpc.ClientConn
}

func NewBookingServiceClient() (*BookingServiceClient, error) {
	addr := env.GetString("BOOKING_SERVICE_ADDRESS", "booking-service:50052")

	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(addr, dialOptions...)
	if err != nil {
		return nil, err
	}

	client := bookingpb.NewBookingServiceClient(conn)

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
