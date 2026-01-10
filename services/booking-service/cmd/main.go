package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/baobei23/e-ticket/services/booking-service/internal/infrastructure/clients"
	"github.com/baobei23/e-ticket/services/booking-service/internal/infrastructure/events"
	"github.com/baobei23/e-ticket/services/booking-service/internal/infrastructure/grpc"
	"github.com/baobei23/e-ticket/services/booking-service/internal/infrastructure/repository"
	"github.com/baobei23/e-ticket/services/booking-service/internal/service"
	"github.com/baobei23/e-ticket/shared/contracts"
	"github.com/baobei23/e-ticket/shared/messaging"

	grpcserver "google.golang.org/grpc"
)

func main() {
	// 1. Init Listener
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	eventAdapter, err := clients.NewEventGRPCClient()
	if err != nil {
		log.Fatalf("failed to init event client: %v", err)
	}

	// 1. Init RabbitMQ
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://admin:admin@rabbitmq:5672/"
	}
	mqClient, err := messaging.NewRabbitMQClient(amqpURL)
	if err != nil {
		log.Fatalf("Failed to init RabbitMQ: %v", err)
	}
	defer mqClient.Close()

	// Declare Queue agar aman saat publish pertama kali
	mqClient.DeclareQueue(contracts.QueueBookingCreated)

	// 2. Init Publisher
	publisher := events.NewBookingEventPublisher(mqClient)

	// 3. Init Service
	repo := repository.NewInMemoryBookingRepository()
	svc := service.NewBookingService(repo, eventAdapter, publisher)

	// 4. Init gRPC Server
	grpcServer := grpcserver.NewServer()

	// Register Handler (Self-registering pattern)
	grpc.NewBookingHandler(grpcServer, svc)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Booking Service listening on :50052\n")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-stop
	log.Println("\nShutting down Booking Service...")

	grpcServer.GracefulStop()
	log.Println("Booking Service exited properly")
}
