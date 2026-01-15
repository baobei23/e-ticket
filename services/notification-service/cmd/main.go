package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/e-ticket/services/notification-service/internal/infrastructure/events"
	"github.com/baobei23/e-ticket/services/notification-service/internal/infrastructure/mailer"
	"github.com/baobei23/e-ticket/services/notification-service/internal/service"
	"github.com/baobei23/e-ticket/shared/env"
	"github.com/baobei23/e-ticket/shared/messaging"
)

func main() {
	// RabbitMQ
	amqpURL := env.GetString("RABBITMQ_URI", "amqp://admin:admin@rabbitmq:5672/")
	mqClient, err := messaging.NewRabbitMQClient(amqpURL)
	if err != nil {
		log.Fatalf("Failed to init RabbitMQ: %v", err)
	}
	defer mqClient.Close()

	// Mail config
	smtpHost := env.GetString("SMTP_HOST", "")
	smtpPort := env.GetString("SMTP_PORT", "587")
	smtpUser := env.GetString("SMTP_USER", "")
	smtpPass := env.GetString("SMTP_PASS", "")
	smtpFrom := env.GetString("SMTP_FROM", "no-reply@example.com")
	activationBaseURL := env.GetString("ACTIVATION_URL_BASE", "http://localhost:8080/activate")
	retryMax := env.GetInt("EMAIL_RETRY_MAX", 3)
	retryBaseMs := env.GetInt("EMAIL_RETRY_BASE_MS", 500)

	emailSender := mailer.NewSMTPMailer(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom)
	notifSvc := service.NewNotificationService(
		emailSender,
		activationBaseURL,
		retryMax,
		time.Duration(retryBaseMs)*time.Millisecond,
	)

	consumer := events.NewUserActivationConsumer(mqClient, notifSvc)
	consumer.Start()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down Notification Service...")
}
