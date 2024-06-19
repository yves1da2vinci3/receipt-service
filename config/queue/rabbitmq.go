package queue

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQConfig() (*RabbitMQConfig, error) {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://admin:admin@localhost:5672/" // default value for development
	}

	connection, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, err
	}

	return &RabbitMQConfig{
		Connection: connection,
		Channel:    channel,
	}, nil
}

func (r *RabbitMQConfig) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
}

func (r *RabbitMQConfig) DeclareQueue(queueName string) (amqp.Queue, error) {
	args := amqp.Table{
		"x-message-ttl": int32(172800000), // 2 days in milliseconds
	}
	queue, err := r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)
	return queue, err
}
