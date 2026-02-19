package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAdapter "github.com/P4ST4S/manga-translator/backend-api/internal/adapters/http"
	"github.com/P4ST4S/manga-translator/backend-api/internal/adapters/queue/asynq"
	"github.com/P4ST4S/manga-translator/backend-api/internal/adapters/repository/postgres"
	"github.com/P4ST4S/manga-translator/backend-api/internal/adapters/worker/python"
	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/config"
	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/database"
	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/logger"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Parse command-line flags
	mode := flag.String("mode", "api", "Run mode: api (HTTP server) or worker (job processor)")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize logger
	zapLogger, err := logger.New(&cfg.Logging)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("starting manga translator",
		zap.String("mode", *mode),
		zap.String("env", cfg.Server.Env),
	)

	// Initialize database
	db, err := database.NewPostgresPool(&cfg.Database, zapLogger)
	if err != nil {
		zapLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	requestRepo := postgres.NewRequestRepository(db)
	resultRepo := postgres.NewResultRepository(db)

	// Initialize queue client
	queueClient, err := asynq.NewQueueClient(&cfg.Redis, zapLogger)
	if err != nil {
		zapLogger.Fatal("failed to initialize queue client", zap.Error(err))
	}
	defer queueClient.Close()

	// Run in selected mode
	switch *mode {
	case "worker":
		runWorker(cfg, zapLogger, requestRepo, resultRepo)
	case "api":
		fallthrough
	default:
		runAPI(cfg, zapLogger, requestRepo, resultRepo, queueClient)
	}
}

func runAPI(
	cfg *config.Config,
	logger *zap.Logger,
	requestRepo ports.RequestRepository,
	resultRepo ports.ResultRepository,
	queueClient ports.QueueClient,
) {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Manga Translator API",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		BodyLimit:    int(cfg.Storage.MaxUploadSize),
		ErrorHandler: customErrorHandler(logger),
	})

	// Setup routes
	httpAdapter.SetupRoutes(app, cfg, logger, requestRepo, resultRepo, queueClient)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		logger.Info("API server listening", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down API server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("server shutdown failed", zap.Error(err))
	}

	logger.Info("API server stopped")
}

func runWorker(
	cfg *config.Config,
	logger *zap.Logger,
	requestRepo ports.RequestRepository,
	resultRepo ports.ResultRepository,
) {
	// Initialize Python executor
	executor := python.NewPythonExecutor(&cfg.Worker, &cfg.Storage, logger)

	// Initialize queue server
	queueServer := asynq.NewQueueServer(cfg, logger, requestRepo, resultRepo, executor)

	// Start worker in goroutine
	go func() {
		logger.Info("starting asynq worker")
		if err := queueServer.Start(); err != nil {
			logger.Fatal("worker failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down worker...")
	queueServer.Stop()
	logger.Info("worker stopped")
}

func customErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		logger.Error("request error",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.Int("status", code),
		)

		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
