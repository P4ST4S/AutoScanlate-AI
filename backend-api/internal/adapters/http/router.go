package http

import (
	"github.com/P4ST4S/manga-translator/backend-api/internal/adapters/http/handlers"
	"github.com/P4ST4S/manga-translator/backend-api/internal/adapters/http/middleware"
	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/config"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// SetupRoutes sets up all HTTP routes
func SetupRoutes(
	app *fiber.App,
	cfg *config.Config,
	logger *zap.Logger,
	requestRepo ports.RequestRepository,
	resultRepo ports.ResultRepository,
	queueClient ports.QueueClient,
) {
	// Middleware
	app.Use(middleware.Recovery())
	app.Use(middleware.Logger(logger))
	app.Use(middleware.CORS(cfg.CORS.Origins))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API routes
	api := app.Group("/api")

	// Upload handler
	uploadHandler := handlers.NewUploadHandler(requestRepo, queueClient, cfg, logger)
	api.Post("/translate", uploadHandler.Upload)

	// Requests handler
	requestsHandler := handlers.NewRequestsHandler(requestRepo, logger)
	api.Get("/requests", requestsHandler.List)
	api.Get("/requests/:id", requestsHandler.GetByID)

	// Events handler (SSE)
	eventsHandler := handlers.NewEventsHandler(requestRepo, cfg, logger)
	api.Get("/requests/:id/events", eventsHandler.StreamProgress)

	// Results handler
	resultsHandler := handlers.NewResultsHandler(requestRepo, resultRepo, logger)
	api.Get("/results/:id", resultsHandler.GetByRequestID)

	// File serving
	filesHandler := handlers.NewFilesHandler(cfg, logger)
	api.Get("/files/:requestId/:type/*", filesHandler.ServeFile)
}
