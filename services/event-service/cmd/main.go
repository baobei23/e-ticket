package main

import (
	"fmt"
	"log"
	"net"

	"github.com/baobei23/e-ticket/services/event-service/internal/infrastructure/repository"
	"github.com/baobei23/e-ticket/services/event-service/internal/service"
	pb "github.com/baobei23/e-ticket/shared/proto/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Init Listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. Init Dependencies
	repo := repository.NewInMemoryRepository()
	svc := service.NewEventService(repo)

	// 3. Init gRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterEventServiceServer(grpcServer, svc)
	reflection.Register(grpcServer)

	fmt.Printf("Event Service listening on :50051\n")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
