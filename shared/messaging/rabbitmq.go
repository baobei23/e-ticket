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

const (
	ExchangeName       = "e_ticket_exchange"
	ExchangeType       = "topic"
	DeadLetterExchange = "e_ticket_dlx"
	DeadLetterQueue    = "e_ticket_dlq"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	// Retry connection logic
	var conn *amqp.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	client := &RabbitMQClient{conn: conn, ch: ch}

	if err := client.setupTopology(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *RabbitMQClient) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

// Setup Topology: Exchanges, Queues, Bindings
func (c *RabbitMQClient) setupTopology() error {

	err := c.ch.ExchangeDeclare(ExchangeName, ExchangeType, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	c.ch.ExchangeDeclare(DeadLetterExchange, "fanout", true, false, false, false, nil)
	c.ch.QueueDeclare(DeadLetterQueue, true, false, false, false, nil)
	c.ch.QueueBind(DeadLetterQueue, "", DeadLetterExchange, false, nil)

	if err := c.declareAndBind(contracts.QueueBookingCreated, []string{"BookingCreated"}); err != nil {
		return err
	}

	return nil
}

func (c *RabbitMQClient) declareAndBind(queueName string, routingKeys []string) error {
	args := amqp.Table{"x-dead-letter-exchange": DeadLetterExchange}

	q, err := c.ch.QueueDeclare(queueName, true, false, false, false, args)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	for _, key := range routingKeys {
		if err := c.ch.QueueBind(q.Name, key, ExchangeName, false, nil); err != nil {
			return fmt.Errorf("failed to bind queue %s to key %s: %w", queueName, key, err)
		}
	}
	return nil
}

func (c *RabbitMQClient) Publish(ctx context.Context, eventName string, payload interface{}) error {
	envelope := contracts.AmqpMessage{
		EventName: eventName,
		Timestamp: time.Now(),
	}

	payloadBytes, _ := json.Marshal(payload)
	envelope.Payload = payloadBytes

	body, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.ch.PublishWithContext(ctx,
		ExchangeName,
		eventName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *RabbitMQClient) Consume(queueName string, handler func(contracts.AmqpMessage) error) error {
	msgs, err := c.ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var envelope contracts.AmqpMessage
			if err := json.Unmarshal(d.Body, &envelope); err != nil {
				log.Printf("Error unmarshal: %v", err)
				d.Nack(false, false)
				continue
			}

			if err := handler(envelope); err != nil {
				log.Printf("Error handler: %v", err)
				d.Nack(false, false)
			} else {
				d.Ack(false)
			}
		}
	}()
	return nil
}
