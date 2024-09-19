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
	topics      map[string]*pubsub.Topic
	mutex       sync.Mutex
	projectID   string
	isConnected bool
}

// NewPubSub creates a new PubSub
func NewPubSub(cfg interface{}) (*PubSub, error) {
	cfgPB := cfg.(*config.PubSubConfig)

	ps := &PubSub{
		projectID: cfgPB.ProjectID,
		topics:    make(map[string]*pubsub.Topic),
	}

	if err := ps.connect(); err != nil {
		return nil, fmt.Errorf("error initializing PubSub: %v", err)
	}
	for _, topicID := range cfgPB.Topics {
		if err := ps.connectToTopic(topicID); err != nil {
			return nil, fmt.Errorf("error initializing PubSub: %v", err)
		}
	}
	fmt.Println("Connected to all declared topics")

	return ps, nil
}

func (ps *PubSub) connect() error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if ps.isConnected {
		return nil
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, ps.projectID)
	if err != nil {
		return err
	}

	ps.client = client
	ps.isConnected = true

	return nil
}

func (ps *PubSub) connectToTopic(topicID string) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if ps.client == nil {
		return fmt.Errorf("tried to connect to a topic without a PubSub client")
	}

	topic := ps.client.Topic(topicID)
	if topic == nil {
		return fmt.Errorf("could not find topic %s", topicID)
	}

	ps.topics[topicID] = topic

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

	if !ps.isConnected {
		return fmt.Errorf("tried to publish to a PubSub that is not connected")
	}

	topic, exists := ps.topics[string(queue)]
	if !exists {
		return fmt.Errorf("topic %s not found in configuration", queue)
	}

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
	log.Printf("Published a message to topic %s; msg ID: %v\n", queue, id)
	return nil
}
