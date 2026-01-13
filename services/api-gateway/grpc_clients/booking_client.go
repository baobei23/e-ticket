package grpc_clients

import (
	"github.com/baobei23/e-ticket/shared/env"
	bookingpb "github.com/baobei23/e-ticket/shared/proto/booking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookingServiceClient struct {
	Client bookingpb.BookingServiceClient
	conn   *grpc.ClientConn
}

func NewBookingServiceClient() (*BookingServiceClient, error) {
	svcAddr := env.GetString("BOOKING_SERVICE_ADDRESS", "booking-service:50052")

	conn, err := grpc.NewClient(svcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
