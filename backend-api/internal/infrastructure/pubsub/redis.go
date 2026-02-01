package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// ProgressUpdate represents a progress update message
type ProgressUpdate struct {
	RequestID uuid.UUID `json:"requestId"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	Message   string    `json:"message"`
}

// Publisher handles publishing messages to Redis
type Publisher struct {
	client *redis.Client
	logger *zap.Logger
}

// NewPublisher creates a new Redis publisher
func NewPublisher(cfg *config.RedisConfig, logger *zap.Logger) (*Publisher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis publisher initialized", zap.String("addr", cfg.Addr))

	return &Publisher{
		client: client,
		logger: logger,
	}, nil
}

// PublishProgress publishes a progress update to a request-specific channel
func (p *Publisher) PublishProgress(ctx context.Context, update ProgressUpdate) error {
	channel := fmt.Sprintf("request:%s:progress", update.RequestID)

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal progress update: %w", err)
	}

	if err := p.client.Publish(ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish to Redis: %w", err)
	}

	p.logger.Debug("published progress update",
		zap.String("request_id", update.RequestID.String()),
		zap.Int("progress", update.Progress),
		zap.String("channel", channel),
	)

	return nil
}

// Close closes the Redis client
func (p *Publisher) Close() error {
	return p.client.Close()
}

// Subscriber handles subscribing to Redis channels
type Subscriber struct {
	client *redis.Client
	logger *zap.Logger
}

// NewSubscriber creates a new Redis subscriber
func NewSubscriber(cfg *config.RedisConfig, logger *zap.Logger) (*Subscriber, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis subscriber initialized", zap.String("addr", cfg.Addr))

	return &Subscriber{
		client: client,
		logger: logger,
	}, nil
}

// SubscribeToProgress subscribes to progress updates for a specific request
func (s *Subscriber) SubscribeToProgress(ctx context.Context, requestID uuid.UUID) (<-chan ProgressUpdate, error) {
	channel := fmt.Sprintf("request:%s:progress", requestID)

	pubsub := s.client.Subscribe(ctx, channel)

	// Wait for subscription confirmation
	if _, err := pubsub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	s.logger.Info("subscribed to progress channel",
		zap.String("request_id", requestID.String()),
		zap.String("channel", channel),
	)

	// Create channel for progress updates
	updates := make(chan ProgressUpdate, 10)

	// Start goroutine to receive messages
	go func() {
		defer close(updates)
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				s.logger.Debug("subscription cancelled",
					zap.String("request_id", requestID.String()),
				)
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}

				var update ProgressUpdate
				if err := json.Unmarshal([]byte(msg.Payload), &update); err != nil {
					s.logger.Error("failed to unmarshal progress update",
						zap.Error(err),
						zap.String("payload", msg.Payload),
					)
					continue
				}

				select {
				case updates <- update:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return updates, nil
}

// Close closes the Redis client
func (s *Subscriber) Close() error {
	return s.client.Close()
}
