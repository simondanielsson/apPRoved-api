package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/simondanielsson/apPRoved/cmd/config"
)

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	connStr      string
	mutex        sync.Mutex
	connected    bool
	reconnectDur time.Duration
}

// NewRabbitMQ creates a new RabbitMQ
func NewRabbitMQ(cfg interface{}) (*RabbitMQ, error) {
	cfgMQ := cfg.(*config.RabbitMQConfig)
	queue := &RabbitMQ{
		connStr:      cfgMQ.Url,
		reconnectDur: 5 * time.Second,
	}
	if err := queue.connect(); err != nil {
		log.Fatalf("error initializing RabbitMQ: %v", err)
	}
	queue.declareQueues(cfgMQ)

	return queue, nil
}

// InitRabbitMQ initializes the RabbitMQ connection and channel
func (r *RabbitMQ) connect() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	fmt.Println("Initializing RabbitMQ connection...")

	if r.connected {
		return nil
	}

	conn, err := amqp.Dial(r.connStr)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		return err
	}

	r.conn = conn
	r.channel = channel
	r.connected = true

	fmt.Println("RabbitMQ connection initialized")
	return nil
}

func (r *RabbitMQ) reconnect() {
	for {
		if err := r.connect(); err != nil {
			log.Printf("Failed to reconnect to RabbitMQ: %v. Retrying in %v...\n", err, r.reconnectDur)
			time.Sleep(r.reconnectDur)
		} else {
			log.Println("Successfully reconnected to RabbitMQ.")
			break
		}
	}
}

func (r *RabbitMQ) declareQueues(cfg *config.RabbitMQConfig) error {
	for _, queue := range cfg.Queues {
		_, err := r.channel.QueueDeclare(
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

func (r *RabbitMQ) Close() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	fmt.Println("Closing RabbitMQ connection...")

	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	r.connected = false
}

func (r *RabbitMQ) Publish(ctx context.Context, queue config.QueueName, message interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.connected {
		return amqp.ErrClosed
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	fmt.Printf("Publishing message to queue %s\n", queue)
	err = r.channel.Publish(
		"", // exchange
		string(queue),
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}
