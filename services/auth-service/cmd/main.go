package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/e-ticket/services/auth-service/internal/infrastructure/events"
	"github.com/baobei23/e-ticket/services/auth-service/internal/infrastructure/grpc"
	"github.com/baobei23/e-ticket/services/auth-service/internal/infrastructure/repository"
	"github.com/baobei23/e-ticket/services/auth-service/internal/service"
	"github.com/baobei23/e-ticket/shared/db"
	"github.com/baobei23/e-ticket/shared/env"
	"github.com/baobei23/e-ticket/shared/messaging"
	"github.com/baobei23/e-ticket/shared/tracing"

	grpcserver "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	tracerCfg := tracing.Config{
		ServiceName:    "auth-service",
		Environment:    env.GetString("ENVIRONMENT", "development"),
		JaegerEndpoint: env.GetString("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
	}
	sh, err := tracing.InitTracer(tracerCfg)
	if err != nil {
		log.Fatalf("Failed to initialize the tracer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer sh(ctx)

	// Init Listener
	lis, err := net.Listen("tcp", ":50054") // Port 50054
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Init RabbitMQ
	amqpURL := env.GetString("RABBITMQ_URI", "amqp://admin:admin@rabbitmq:5672/")
	mqClient, err := messaging.NewRabbitMQClient(amqpURL)
	if err != nil {
		log.Fatalf("Failed to init RabbitMQ: %v", err)
	}
	defer mqClient.Close()

	// Init DB
	dbURI := env.GetString("POSTGRES_URI", "postgresql://postgres:postgres@eticket-postgres:5432/auth_service")
	dbPool, err := db.New(dbURI, 10, 5, 10*time.Second, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbPool.Close()

	// Init Dependencies
	repo := repository.NewPostgresRepository(dbPool)
	publisher := events.NewUserActivationPublisher(mqClient)
	svc := service.NewAuthService(repo, publisher)

	// Init gRPC Server
	server := grpcserver.NewServer(tracing.WithTracingInterceptors()...)
	grpc.NewAuthHandler(server, svc)
	reflection.Register(server)

	// Run server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Auth Service listening on :50054\n")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down Auth Service...")
	server.GracefulStop()
	log.Println("Auth Service exited properly")
}
