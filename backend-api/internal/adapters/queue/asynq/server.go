package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/P4ST4S/manga-translator/backend-api/internal/domain"
	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/config"
	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/pubsub"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type queueServer struct {
	server      *asynq.Server
	mux         *asynq.ServeMux
	logger      *zap.Logger
	requestRepo ports.RequestRepository
	resultRepo  ports.ResultRepository
	executor    ports.WorkerExecutor
	storagePath string
	publisher   *pubsub.Publisher
}

// NewQueueServer creates a new Asynq queue server
func NewQueueServer(
	cfg *config.Config,
	logger *zap.Logger,
	requestRepo ports.RequestRepository,
	resultRepo ports.ResultRepository,
	executor ports.WorkerExecutor,
) ports.QueueServer {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: cfg.Worker.Concurrency,
			Queues: map[string]int{
				QueueCritical: 6,
				QueueDefault:  3,
			},
			Logger: &asynqLogger{logger: logger},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.Error("task execution failed",
					zap.String("type", task.Type()),
					zap.Error(err),
				)
			}),
		},
	)

	// Create Redis publisher for progress updates
	publisher, err := pubsub.NewPublisher(&cfg.Redis, logger)
	if err != nil {
		logger.Fatal("failed to create publisher", zap.Error(err))
	}

	qs := &queueServer{
		server:      server,
		mux:         asynq.NewServeMux(),
		logger:      logger,
		requestRepo: requestRepo,
		resultRepo:  resultRepo,
		executor:    executor,
		storagePath: cfg.Storage.Path,
		publisher:   publisher,
	}

	// Register task handlers
	qs.mux.HandleFunc(TaskTypeTranslation, qs.handleTranslationTask)

	logger.Info("asynq server initialized",
		zap.Int("concurrency", cfg.Worker.Concurrency),
	)

	return qs
}

func (qs *queueServer) Start() error {
	qs.logger.Info("starting asynq server")
	return qs.server.Run(qs.mux)
}

func (qs *queueServer) Stop() error {
	qs.logger.Info("stopping asynq server")
	qs.server.Shutdown()
	return nil
}

// handleTranslationTask processes a translation task
func (qs *queueServer) handleTranslationTask(ctx context.Context, task *asynq.Task) error {
	// Parse payload
	var payload TranslationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	requestID := payload.RequestID
	qs.logger.Info("processing translation task",
		zap.String("request_id", requestID.String()),
		zap.String("file_path", payload.FilePath),
		zap.String("file_type", payload.FileType),
	)

	// Update request status to processing
	if err := qs.requestRepo.UpdateStatus(ctx, requestID, domain.StatusProcessing, 0); err != nil {
		qs.logger.Error("failed to update status to processing", zap.Error(err))
		// Continue anyway
	}

	// Execute translation with progress callback
	progressCallback := func(progress int, message string) {
		qs.logger.Info("translation progress",
			zap.String("request_id", requestID.String()),
			zap.Int("progress", progress),
			zap.String("message", message),
		)

		// Update database with progress
		if progress >= 0 {
			if err := qs.requestRepo.UpdateStatus(ctx, requestID, domain.StatusProcessing, progress); err != nil {
				qs.logger.Error("failed to update progress", zap.Error(err))
			}
		}

		// Publish to Redis pub/sub for SSE
		update := pubsub.ProgressUpdate{
			RequestID: requestID,
			Status:    string(domain.StatusProcessing),
			Progress:  progress,
			Message:   message,
		}
		if err := qs.publisher.PublishProgress(ctx, update); err != nil {
			qs.logger.Error("failed to publish progress", zap.Error(err))
		} else {
			qs.logger.Debug("published progress to Redis",
				zap.String("request_id", requestID.String()),
				zap.Int("progress", progress),
			)
		}
	}

	output, err := qs.executor.Translate(ctx, payload.FilePath, progressCallback)
	if err != nil {
		qs.logger.Error("translation failed",
			zap.String("request_id", requestID.String()),
			zap.Error(err),
		)

		// Update request with error
		req, _ := qs.requestRepo.GetByID(ctx, requestID)
		if req != nil {
			req.SetError(err.Error())
			qs.requestRepo.Update(ctx, req)
		}

		// Publish error event
		errorUpdate := pubsub.ProgressUpdate{
			RequestID: requestID,
			Status:    string(domain.StatusFailed),
			Progress:  0,
			Message:   fmt.Sprintf("Translation failed: %s", err.Error()),
		}
		if pubErr := qs.publisher.PublishProgress(ctx, errorUpdate); pubErr != nil {
			qs.logger.Error("failed to publish error", zap.Error(pubErr))
		}

		return fmt.Errorf("translation failed: %w", err)
	}

	// Process output files
	if err := qs.processOutputFiles(ctx, requestID, payload.FileType, output); err != nil {
		qs.logger.Error("failed to process output files",
			zap.String("request_id", requestID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to process output: %w", err)
	}

	// Update request to completed
	if err := qs.requestRepo.UpdateStatus(ctx, requestID, domain.StatusCompleted, 100); err != nil {
		qs.logger.Error("failed to update status to completed", zap.Error(err))
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Publish completion event
	completeUpdate := pubsub.ProgressUpdate{
		RequestID: requestID,
		Status:    string(domain.StatusCompleted),
		Progress:  100,
		Message:   "Translation completed successfully",
	}
	if err := qs.publisher.PublishProgress(ctx, completeUpdate); err != nil {
		qs.logger.Error("failed to publish completion", zap.Error(err))
	}

	qs.logger.Info("translation task completed",
		zap.String("request_id", requestID.String()),
	)

	return nil
}

func (qs *queueServer) processOutputFiles(
	ctx context.Context,
	requestID uuid.UUID,
	fileType string,
	output *ports.TranslationOutput,
) error {
	// Create output directories
	originalsDir := filepath.Join(qs.storagePath, "originals", requestID.String())
	translatedDir := filepath.Join(qs.storagePath, "translated", requestID.String())

	for _, dir := range []string{originalsDir, translatedDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// For ZIP files, extract and organize pages
	if fileType == "zip" {
		// TODO: Extract ZIP and catalog pages
		// For now, just move the translated ZIP
		destPath := filepath.Join(translatedDir, filepath.Base(output.OutputPath))
		if err := copyFile(output.OutputPath, destPath); err != nil {
			return fmt.Errorf("failed to copy translated zip: %w", err)
		}

		qs.logger.Info("processed zip output",
			zap.String("request_id", requestID.String()),
			zap.String("output_path", destPath),
		)

		return nil
	}

	// For single images, create result entries
	results := make([]*domain.Result, 0, len(output.Pages))
	for _, page := range output.Pages {
		// Copy files to storage
		originalDest := filepath.Join(originalsDir, filepath.Base(page.OriginalPath))
		translatedDest := filepath.Join(translatedDir, filepath.Base(page.TranslatedPath))

		if err := copyFile(page.OriginalPath, originalDest); err != nil {
			return fmt.Errorf("failed to copy original: %w", err)
		}

		if err := copyFile(page.TranslatedPath, translatedDest); err != nil {
			return fmt.Errorf("failed to copy translated: %w", err)
		}

		// Create API paths
		originalAPIPath := fmt.Sprintf("/api/files/%s/originals/%s", requestID, filepath.Base(originalDest))
		translatedAPIPath := fmt.Sprintf("/api/files/%s/translated/%s", requestID, filepath.Base(translatedDest))

		result := domain.NewResult(requestID, page.PageNumber, originalAPIPath, translatedAPIPath)
		results = append(results, result)
	}

	// Save results to database
	if len(results) > 0 {
		if err := qs.resultRepo.CreateBatch(ctx, results); err != nil {
			return fmt.Errorf("failed to save results: %w", err)
		}

		// Update page count
		req, err := qs.requestRepo.GetByID(ctx, requestID)
		if err == nil {
			req.PageCount = len(results)
			qs.requestRepo.Update(ctx, req)
		}
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// asynqLogger adapts zap.Logger to asynq.Logger interface
type asynqLogger struct {
	logger *zap.Logger
}

func (l *asynqLogger) Debug(args ...interface{}) {
	l.logger.Sugar().Debug(args...)
}

func (l *asynqLogger) Info(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}

func (l *asynqLogger) Warn(args ...interface{}) {
	l.logger.Sugar().Warn(args...)
}

func (l *asynqLogger) Error(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

func (l *asynqLogger) Fatal(args ...interface{}) {
	l.logger.Sugar().Fatal(args...)
}
