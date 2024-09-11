package mq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/simondanielsson/apPRoved/cmd/config"
)

var (
	rabbitMQConn    *amqp.Connection
	rabbitMQChannel *amqp.Channel
)

// InitRabbitMQ initializes the RabbitMQ connection and channel
func InitRabbitMQ(cfg *config.RabbitMQConfig) error {
	var err error

	rabbitMQConn, err = amqp.Dial(cfg.Url)
	if err != nil {
		return err
	}

	rabbitMQChannel, err = rabbitMQConn.Channel()
	if err != nil {
		return err
	}

	for _, queue := range cfg.Queues {
		_, err = rabbitMQChannel.QueueDeclare(
			queue.Name,
			queue.Durable,
			queue.AutoDelete,
			queue.Exclusive,
			queue.NoWait,
			queue.Args,
		)
		if err != nil {
			return err
		}
		log.Printf("Queue declared: %s", queue.Name)
	}

	return nil
}

func CloseRabbitMQ() {
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}

	if rabbitMQConn != nil {
		rabbitMQConn.Close()
	}
}

func Publish(queue config.QueueName, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = rabbitMQChannel.Publish(
		"", // exchange
		string(queue),
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Sent message on queue %s", queue)
	return nil
}
