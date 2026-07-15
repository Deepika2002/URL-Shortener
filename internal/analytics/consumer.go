package analytics

import (
	"context"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"urlshortener/pkg/config"
)

// Repository defines the data sink for the analytics consumer (e.g., ClickHouse).
type Repository interface {
	RecordEvent(ctx context.Context, eventType string, payload []byte) error
}

// Consumer defines the background worker for processing analytics events.
type Consumer interface {
	Start(ctx context.Context) error
	Close() error
}

type consumerImpl struct {
	client sarama.ConsumerGroup
	repo   Repository
}

// NewConsumer creates a new Kafka/Redpanda consumer group for analytics.
func NewConsumer(cfg *config.Config, repo Repository) (Consumer, error) {
	brokers := strings.Split(cfg.KafkaBrokers, ",")

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0 // Set to a standard recent Kafka version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, "analytics-group", config)
	if err != nil {
		return nil, err
	}

	return &consumerImpl{
		client: client,
		repo:   repo,
	}, nil
}

// Start begins the blocking consumer loop. It should be run in a separate goroutine.
func (c *consumerImpl) Start(ctx context.Context) error {
	topics := []string{"url_created", "url_redirected"}
	handler := &consumerGroupHandler{repo: c.repo}

	log.Printf("Analytics Consumer Group starting on topics: %v\n", topics)

	for {
		// Consume blocks until an error occurs or context is canceled.
		// It creates a new session and re-joins the consumer group automatically if needed.
		if err := c.client.Consume(ctx, topics, handler); err != nil {
			if err == sarama.ErrClosedConsumerGroup {
				return nil
			}
			log.Printf("Error from consumer: %v\n", err)
			return err
		}

		// Check if context was cancelled, signaling that the consumer should completely stop.
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Close gracefully stops the consumer group.
func (c *consumerImpl) Close() error {
	return c.client.Close()
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	repo Repository
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from the Kafka topic(s)
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Println("Message channel was closed")
				return nil
			}

			// Delegate the raw JSON bytes to the repository based on the topic
			if h.repo != nil {
				if message.Topic == "url_created" {
					_ = h.repo.RecordEvent(session.Context(), "created", message.Value)
				} else if message.Topic == "url_redirected" {
					_ = h.repo.RecordEvent(session.Context(), "redirected", message.Value)
				}
			}

			// Mark the message as successfully processed so the offset is committed
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			// Session was cancelled/rebalanced
			return nil
		}
	}
}
