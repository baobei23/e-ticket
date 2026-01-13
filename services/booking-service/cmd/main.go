package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/e-ticket/services/booking-service/internal/infrastructure/clients"
	"github.com/baobei23/e-ticket/services/booking-service/internal/infrastructure/events"
	"github.com/baobei23/e-ticket/services/booking-service/internal/infrastructure/grpc"
	"github.com/baobei23/e-ticket/services/booking-service/internal/infrastructure/repository"
	"github.com/baobei23/e-ticket/services/booking-service/internal/service"
	"github.com/baobei23/e-ticket/shared/db"
	"github.com/baobei23/e-ticket/shared/env"
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

	paymentAdapter, err := clients.NewPaymentGRPCClient()
	if err != nil {
		log.Fatalf("failed to init payment client: %v", err)
	}

	//Init RabbitMQ
	amqpURL := env.GetString("RABBITMQ_URI", "amqp://admin:admin@rabbitmq:5672/")
	mqClient, err := messaging.NewRabbitMQClient(amqpURL)
	if err != nil {
		log.Fatalf("Failed to init RabbitMQ: %v", err)
	}
	defer mqClient.Close()

	dbURI := env.GetString("POSTGRES_URI", "postgresql://postgres:postgres@eticket-postgres:5432/booking_service")
	log.Println("Connecting to database...")
	pool, err := db.New(dbURI, 10, 5, 10*time.Second, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Init Publisher
	publisher := events.NewBookingEventPublisher(mqClient)

	// Init Service
	repo := repository.NewPostgresRepository(pool)
	svc := service.NewBookingService(repo, eventAdapter, publisher, paymentAdapter)

	paymentConsumer := events.NewPaymentConsumer(mqClient, svc)
	paymentConsumer.Start()

	// Init gRPC Server
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
