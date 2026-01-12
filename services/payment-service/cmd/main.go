package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/baobei23/e-ticket/services/payment-service/internal/infrastructure/events"
	"github.com/baobei23/e-ticket/services/payment-service/internal/infrastructure/gateway"
	"github.com/baobei23/e-ticket/services/payment-service/internal/infrastructure/grpc"
	"github.com/baobei23/e-ticket/services/payment-service/internal/infrastructure/repository"
	"github.com/baobei23/e-ticket/services/payment-service/internal/service"
	"github.com/baobei23/e-ticket/shared/env"
	"github.com/baobei23/e-ticket/shared/messaging"

	grpcserver "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	amqpURL := env.GetString("RABBITMQ_URI", "amqp://admin:admin@rabbitmq:5672/")
	mqClient, err := messaging.NewRabbitMQClient(amqpURL)
	if err != nil {
		log.Fatalf("Failed to init RabbitMQ: %v", err)
	}
	defer mqClient.Close()
	publisher := events.NewPaymentEventPublisher(mqClient)

	// 1. Init Listener (Port 50053)
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. Setup Stripe Gateway
	stripeKey := env.GetString("STRIPE_API_KEY", "")
	successURL := env.GetString("STRIPE_SUCCESS_URL", "http://localhost:8080/success")
	cancelURL := env.GetString("STRIPE_CANCEL_URL", "http://localhost:8080/cancel")

	stripeGateway := gateway.NewStripeGateway(stripeKey, successURL, cancelURL)

	// 3. Init Dependencies
	repo := repository.NewInMemoryPaymentRepository()
	svc := service.NewPaymentService(repo, stripeGateway, publisher)

	// 4. Init gRPC Server
	server := grpcserver.NewServer()
	grpc.NewPaymentHandler(server, svc)
	reflection.Register(server)

	// 5. Run Server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Payment Service listening on :50053\n")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-stop
	log.Println("\nShutting down Payment Service...")
	server.GracefulStop()
	log.Println("Payment Service exited properly")
}
