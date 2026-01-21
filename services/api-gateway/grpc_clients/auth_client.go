package grpc_clients

import (
	"github.com/baobei23/e-ticket/shared/env"
	"github.com/baobei23/e-ticket/shared/proto/auth"
	"github.com/baobei23/e-ticket/shared/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	Client auth.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthServiceClient() (*AuthServiceClient, error) {
	addr := env.GetString("AUTH_SERVICE_ADDRESS", "auth-service:50054")

	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	conn, err := grpc.NewClient(addr, dialOptions...)
	if err != nil {
		return nil, err
	}

	client := auth.NewAuthServiceClient(conn)
	return &AuthServiceClient{Client: client, conn: conn}, nil
}

func (c *AuthServiceClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
