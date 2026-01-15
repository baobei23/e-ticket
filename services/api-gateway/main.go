package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/e-ticket/services/api-gateway/grpc_clients"
	"github.com/baobei23/e-ticket/shared/env"
	"github.com/gin-gonic/gin"
)

var httpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8080")

func main() {

	registry, err := grpc_clients.NewServiceRegistry()
	if err != nil {
		log.Fatal(err)
	}
	defer registry.Close()

	grpc := NewGatewayServer(registry.Event, registry.Booking, registry.Payment, registry.Auth)

	r := gin.Default()

	r.GET("/health", healthCheckHandler)
	r.GET("/events", grpc.getEventsHandler)
	r.GET("/events/:id", grpc.getEventDetailHandler)
	r.GET("/events/:id/check", grpc.checkAvailabilityHandler)

	r.POST("/auth/register", grpc.RegisterHandler)
	r.POST("/auth/activate", grpc.ActivateHandler)
	r.POST("/auth/login", grpc.LoginHandler)

	protected := r.Group("/")
	protected.Use(AuthMiddleware(registry.Auth))
	{
		protected.POST("/bookings", grpc.CreateBookingHandler)
		protected.GET("/booking/:id", grpc.GetBookingDetailHandler)
	}

	r.POST("/stripe/webhook", grpc.HandleStripeWebhook)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: r,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("API Gateway listening on %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the API Gateway: %v", err)

	case sig := <-shutdown:
		log.Printf("API Gateway is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the API Gateway gracefully: %v", err)
			server.Close()
		}
	}

}
