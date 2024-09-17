package mq

import (
	"context"
	"fmt"

	"github.com/simondanielsson/apPRoved/cmd/config"
)

type MessageQueue interface {
	Close()
	Publish(ctx context.Context, queue config.QueueName, message interface{}) error
}

func NewMessageQueue(cfg *config.Config) (MessageQueue, error) {
	var err error
	var queue MessageQueue

	switch cfg.Server.AMQPMode {
	case "rabbitmq":
		queue, err = NewRabbitMQ(cfg.MQ)
	case "pubsub":
		queue, err = NewPubSub(cfg.PubSub)
	default:
		return nil, fmt.Errorf("invalid AMQP mode: %s. Expected rabbitmq or pubsub", cfg.Server.AMQPMode)
	}

	if err != nil {
		return nil, fmt.Errorf("error initializing message queue: %v", err)
	}

	return queue, nil
}
