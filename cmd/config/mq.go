package config

import "fmt"

type QueueName string

const (
	QueueFileDiffs QueueName = "review-file-diffs"
)

var ValidQueueNames = map[string]QueueName{
	string(QueueFileDiffs): QueueFileDiffs,
}

// ValidateRabbitMQConfig validates that all queue names in the config are valid
func ValidateRabbitMQConfig(cfg *RabbitMQConfig) error {
	for _, queue := range cfg.Queues {
		if _, exists := ValidQueueNames[queue.Name]; !exists {
			return fmt.Errorf("invalid queue name in configuration: %s", queue.Name)
		}
	}
	return nil
}
