package clients

import (
	"context"
	"fmt"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	"github.com/baobei23/e-ticket/shared/env"
	eventpb "github.com/baobei23/e-ticket/shared/proto/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventGRPCClient struct {
	client eventpb.EventServiceClient
	conn   *grpc.ClientConn
}

func NewEventGRPCClient() (domain.EventProvider, error) {

	addr := env.GetString("EVENT_SERVICE_ADDRESS", "event-service:50051")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to event service: %w", err)
	}

	client := eventpb.NewEventServiceClient(conn)

	return &EventGRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *EventGRPCClient) Close() error {
	return c.conn.Close()
}

func (c *EventGRPCClient) CheckAvailability(ctx context.Context, eventID int64, quantity int32) (bool, int64, error) {
	resp, err := c.client.CheckAvailability(ctx, &eventpb.CheckAvailabilityRequest{
		EventId:  eventID,
		Quantity: quantity,
	})

	if err != nil {
		return false, 0, err
	}

	return resp.IsAvailable, resp.UnitPrice, nil
}

func (c *EventGRPCClient) ReserveSeat(ctx context.Context, eventID int64, quantity int32, bookingID string) (int64, error) {
	resp, err := c.client.ReserveSeat(ctx, &eventpb.ReserveSeatRequest{
		EventId:   eventID,
		Quantity:  quantity,
		BookingId: bookingID,
	})
	if err != nil {
		return 0, err
	}
	if !resp.Reserved {
		return 0, fmt.Errorf("failed to reserve seat: %s", resp.Message)
	}
	return resp.UnitPrice, nil
}

func (c *EventGRPCClient) ReleaseSeat(ctx context.Context, eventID int64, quantity int32, bookingID string) error {
	resp, err := c.client.ReleaseSeat(ctx, &eventpb.ReleaseSeatRequest{
		EventId:   eventID,
		Quantity:  quantity,
		BookingId: bookingID,
	})
	if err != nil {
		return err
	}
	if !resp.Released {
		return fmt.Errorf("failed to release seat: %s", resp.Message)
	}
	return nil
}
