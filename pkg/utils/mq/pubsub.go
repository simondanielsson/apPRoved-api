package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/simondanielsson/apPRoved/cmd/config"
)

type PubSub struct {
	client      *pubsub.Client
	mutex       sync.Mutex
	ProjectID   string
	isConnected bool
}

// NewPubSub creates a new PubSub
func NewPubSub(cfg interface{}) (*PubSub, error) {
	cfgPB := cfg.(*config.PubSubConfig)

	ps := &PubSub{
		ProjectID: cfgPB.ProjectID,
	}

	if err := ps.connect(); err != nil {
		return nil, fmt.Errorf("error initializing PubSub: %v", err)
	}

	return ps, nil
}

func (ps *PubSub) connect() error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, ps.ProjectID)
	if err != nil {
		return err
	}

	ps.client = client
	ps.isConnected = true

	return nil
}

func (ps *PubSub) Close() {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if ps.client != nil {
		ps.client.Close()
	}

	ps.isConnected = false
}

func (ps *PubSub) Publish(ctx context.Context, queue config.QueueName, message interface{}) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	topic := ps.client.Topic(string(queue))
	defer topic.Stop()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageBytes,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	log.Printf("Published a message; msg ID: %v\n", id)
	return nil
}
