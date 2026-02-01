package asynq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/config"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

const (
	// Task types
	TaskTypeTranslation = "translation:process"

	// Queue names
	QueueCritical = "critical"
	QueueDefault  = "default"
)

// TranslationPayload represents the payload for a translation task
type TranslationPayload struct {
	RequestID uuid.UUID `json:"requestId"`
	FilePath  string    `json:"filePath"`
	FileType  string    `json:"fileType"`
}

type queueClient struct {
	client *asynq.Client
	logger *zap.Logger
}

// NewQueueClient creates a new Asynq queue client
func NewQueueClient(cfg *config.RedisConfig, logger *zap.Logger) (ports.QueueClient, error) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	client := asynq.NewClient(redisOpt)

	logger.Info("asynq client initialized",
		zap.String("redis_addr", cfg.Addr),
	)

	return &queueClient{
		client: client,
		logger: logger,
	}, nil
}

func (q *queueClient) EnqueueTranslation(ctx context.Context, requestID uuid.UUID, filePath string, fileType string) error {
	payload := TranslationPayload{
		RequestID: requestID,
		FilePath:  filePath,
		FileType:  fileType,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskTypeTranslation, payloadBytes)

	// Enqueue task with options
	info, err := q.client.EnqueueContext(ctx, task,
		asynq.Queue(QueueDefault),
		asynq.MaxRetry(3),
		asynq.Timeout(0), // No timeout, let worker config handle it
	)

	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	q.logger.Info("translation task enqueued",
		zap.String("request_id", requestID.String()),
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
	)

	return nil
}

func (q *queueClient) Close() error {
	return q.client.Close()
}
