package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/e-ticket/shared/env"
	"github.com/gin-gonic/gin"
)

var httpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8080")

func main() {

	r := gin.Default()

	r.GET("/health", healthCheckHandler)

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
