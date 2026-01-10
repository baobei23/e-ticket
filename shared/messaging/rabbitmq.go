package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/baobei23/e-ticket/shared/contracts"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (c *RabbitMQClient) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *RabbitMQClient) DeclareQueue(name string) error {
	_, err := c.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (c *RabbitMQClient) Publish(ctx context.Context, queueName, eventName string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	envelope := contracts.AmqpMessage{
		EventName: eventName,
		Timestamp: time.Now(),
		Payload:   payloadBytes,
	}

	body, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal envelope: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *RabbitMQClient) Consume(queueName string, handler func(contracts.AmqpMessage) error) error {
	msgs, err := c.ch.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack (FALSE = Manual Ack)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var envelope contracts.AmqpMessage
			if err := json.Unmarshal(d.Body, &envelope); err != nil {
				log.Printf("Error unmarshal envelope: %v", err)
				d.Nack(false, false)
				continue
			}

			if err := handler(envelope); err != nil {
				log.Printf("Error processing message: %v", err)
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}
