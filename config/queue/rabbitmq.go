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
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQConfig{
		Connection: conn,
		Channel:    ch,
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
		"x-dead-letter-exchange": "dlx_exchange",
		"x-message-ttl":          int32(172800000), // 2 days in milliseconds
	}

	queue, err := r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)
	if err != nil {
		return queue, err
	}

	// Declare the dead-letter queue
	_, err = r.Channel.QueueDeclare(
		"dlx_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return queue, err
	}

	// Bind the dead-letter queue to the dead-letter exchange
	err = r.Channel.ExchangeDeclare(
		"dlx_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return queue, err
	}

	err = r.Channel.QueueBind(
		"dlx_queue",
		"",
		"dlx_exchange",
		false,
		nil,
	)
	if err != nil {
		return queue, err
	}

	return queue, nil
}
