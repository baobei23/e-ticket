package clients

import (
	"context"
	"fmt"
	"os"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	pb "github.com/baobei23/e-ticket/shared/proto/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventGRPCClient struct {
	client pb.EventServiceClient
	conn   *grpc.ClientConn
}

func NewEventGRPCClient() (domain.EventProvider, error) {

	addr := os.Getenv("EVENT_SERVICE_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to event service: %w", err)
	}

	client := pb.NewEventServiceClient(conn)

	return &EventGRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *EventGRPCClient) Close() error {
	return c.conn.Close()
}

func (c *EventGRPCClient) CheckAvailability(ctx context.Context, eventID int64, quantity int32) (bool, float64, error) {
	resp, err := c.client.CheckAvailability(ctx, &pb.CheckAvailabilityRequest{
		EventId:  eventID,
		Quantity: quantity,
	})

	if err != nil {
		return false, 0, err
	}

	return resp.IsAvailable, resp.UnitPrice, nil
}
