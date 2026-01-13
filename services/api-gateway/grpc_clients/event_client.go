package grpc_clients

import (
	"github.com/baobei23/e-ticket/shared/env"
	eventpb "github.com/baobei23/e-ticket/shared/proto/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventServiceClient struct {
	Client eventpb.EventServiceClient
	conn   *grpc.ClientConn
}

func NewEventServiceClient() (*EventServiceClient, error) {
	eventServiceAddress := env.GetString("EVENT_SERVICE_ADDRESS", "event-service:50051")

	conn, err := grpc.NewClient(eventServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := eventpb.NewEventServiceClient(conn)

	return &EventServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *EventServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
