package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/e-ticket/services/auth-service/internal/infrastructure/grpc"
	"github.com/baobei23/e-ticket/services/auth-service/internal/infrastructure/repository"
	"github.com/baobei23/e-ticket/services/auth-service/internal/service"
	"github.com/baobei23/e-ticket/shared/db"
	"github.com/baobei23/e-ticket/shared/env"

	grpcserver "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Init Listener
	lis, err := net.Listen("tcp", ":50054") // Port 50054
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Init DB
	dbURI := env.GetString("POSTGRES_URI", "postgresql://postgres:postgres@eticket-postgres:5432/auth_service")
	dbPool, err := db.New(dbURI, 10, 5, 10*time.Second, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbPool.Close()

	// Init Dependencies
	repo := repository.NewPostgresRepository(dbPool)
	svc := service.NewAuthService(repo)

	// Init gRPC Server
	server := grpcserver.NewServer()
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
}
