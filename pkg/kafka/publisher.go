package kafka

import (
	"strings"
	"time"

	"github.com/IBM/sarama"
	"urlshortener/pkg/config"
)

// Publisher defines the interface for publishing messages to Kafka/Redpanda.
type Publisher interface {
	Close() error
	Publish(topic string, key string, message []byte) error
}

type publisher struct {
	producer sarama.SyncProducer
}

// NewPublisher creates and returns a new Publisher connected to Redpanda.
func NewPublisher(cfg *config.Config) (Publisher, error) {
	brokers := strings.Split(cfg.KafkaBrokers, ",")

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Net.DialTimeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &publisher{producer: producer}, nil
}

// Publish sends a message to the specified topic.
func (p *publisher) Publish(topic string, key string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(message),
	}
	_, _, err := p.producer.SendMessage(msg)
	return err
}

// Close gracefully shuts down the producer.
func (p *publisher) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}
