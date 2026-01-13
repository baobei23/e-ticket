package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/e-ticket/services/event-service/internal/infrastructure/events"
	"github.com/baobei23/e-ticket/services/event-service/internal/infrastructure/grpc"
	"github.com/baobei23/e-ticket/services/event-service/internal/infrastructure/repository"
	"github.com/baobei23/e-ticket/services/event-service/internal/service"
	"github.com/baobei23/e-ticket/shared/db"
	"github.com/baobei23/e-ticket/shared/env"
	"github.com/baobei23/e-ticket/shared/messaging"
	grpcserver "google.golang.org/grpc"
)

func main() {

	// Init RabbitMQ
	amqpURL := env.GetString("RABBITMQ_URI", "amqp://admin:admin@rabbitmq:5672/")
	mqClient, err := messaging.NewRabbitMQClient(amqpURL)
	if err != nil {
		log.Fatalf("Failed to init RabbitMQ: %v", err)
	}
	defer mqClient.Close()

	// Init Listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Init Database
	dbURI := env.GetString("POSTGRES_URI", "postgresql://postgres:postgres@eticket-postgres:5432/event_service")
	log.Println("Connecting to database...")
	pool, err := db.New(dbURI, 10, 5, 10*time.Second, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Init Dependencies
	repo := repository.NewPostgresRepository(pool)
	service := service.NewEventService(repo)

	consumer := events.NewEventConsumer(mqClient, service)
	consumer.Start()

	// Init gRPC Server
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
