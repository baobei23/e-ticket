package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/baobei23/e-ticket/services/event-service/internal/infrastructure/grpc"
	"github.com/baobei23/e-ticket/services/event-service/internal/infrastructure/repository"
	"github.com/baobei23/e-ticket/services/event-service/internal/service"
	grpcserver "google.golang.org/grpc"
)

func main() {
	// 1. Init Listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. Init Dependencies
	repo := repository.NewInMemoryRepository()
	service := service.NewEventService(repo)

	// 3. Init gRPC Server
	grpcServer := grpcserver.NewServer()
	grpc.NewEventHandler(grpcServer, service)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Event Service listening on :50051\n")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-stop
	log.Println("\nShutting down Event Service...")

	grpcServer.GracefulStop()
	log.Println("Event Service exited properly")
}
