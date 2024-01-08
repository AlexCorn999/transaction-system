package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	Ch   *amqp.Channel
}

// NewRabbitMQ creates a connection to rabbitmq.
func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial("amqp://restapi:1234@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("unable to open connect to RabbitMQ server. Error: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel. Error: %w", err)
	}

	_, err = ch.QueueDeclare(
		"transactions",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue. Error: %w", err)
	}

	return &RabbitMQ{
		conn: conn,
		Ch:   ch,
	}, nil
}

// Close closing the rabbitmq connection.
func (r *RabbitMQ) Close() error {
	if err := r.Ch.Close(); err != nil {
		return err
	}

	if err := r.conn.Close(); err != nil {
		return err
	}

	return nil
}
