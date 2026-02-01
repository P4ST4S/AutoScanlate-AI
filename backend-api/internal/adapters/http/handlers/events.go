package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/P4ST4S/manga-translator/backend-api/internal/domain"
	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/config"
	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/pubsub"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EventsHandler struct {
	requestRepo ports.RequestRepository
	cfg         *config.Config
	logger      *zap.Logger
}

func NewEventsHandler(
	requestRepo ports.RequestRepository,
	cfg *config.Config,
	logger *zap.Logger,
) *EventsHandler {
	return &EventsHandler{
		requestRepo: requestRepo,
		cfg:         cfg,
		logger:      logger,
	}
}

// StreamProgress handles GET /api/requests/:id/events
// Streams progress updates via Server-Sent Events (SSE)
func (h *EventsHandler) StreamProgress(c *fiber.Ctx) error {
	// Parse request ID
	idStr := c.Params("id")
	requestID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request ID",
		})
	}

	// Verify request exists
	request, err := h.requestRepo.GetByID(c.Context(), requestID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "request not found",
			})
		}
		h.logger.Error("failed to get request", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve request",
		})
	}

	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// If already completed, send final status and close
	if request.IsCompleted() {
		h.sendSSEEvent(c, "complete", fiber.Map{
			"status":   request.Status,
			"progress": request.Progress,
			"message":  fmt.Sprintf("Translation %s", request.Status),
		})
		return nil
	}

	// Create Redis subscriber
	subscriber, err := pubsub.NewSubscriber(&h.cfg.Redis, h.logger)
	if err != nil {
		h.logger.Error("failed to create subscriber", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to initialize event stream",
		})
	}
	// Note: Don't defer close here - it will be closed inside the stream writer

	// Subscribe to progress updates
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Minute)
	// Note: Don't defer cancel here either - it will be called inside the stream writer

	updates, err := subscriber.SubscribeToProgress(ctx, requestID)
	if err != nil {
		h.logger.Error("failed to subscribe", zap.Error(err))
		subscriber.Close()
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to subscribe to updates",
		})
	}

	h.logger.Info("SSE client connected",
		zap.String("request_id", requestID.String()),
	)

	// Setup streaming
	c.Context().Response.Header.Set("X-Accel-Buffering", "no")
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Ensure cleanup when stream ends
		defer subscriber.Close()
		defer cancel()

		// Send initial connected event
		if err := h.writeSSEEvent(w, "connected", fiber.Map{
			"status":   request.Status,
			"progress": request.Progress,
			"message":  "Connected to progress stream",
		}); err != nil {
			h.logger.Error("failed to send connected event", zap.Error(err))
			return
		}

		// Keep-alive ticker
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				h.logger.Info("SSE client disconnected",
					zap.String("request_id", requestID.String()),
				)
				return

			case update, ok := <-updates:
				if !ok {
					// Channel closed
					h.logger.Debug("SSE updates channel closed",
						zap.String("request_id", requestID.String()),
					)
					return
				}

				h.logger.Debug("SSE received update",
					zap.String("request_id", requestID.String()),
					zap.Int("progress", update.Progress),
					zap.String("message", update.Message),
				)

				// Determine event type
				eventType := "progress"
				if update.Status == string(domain.StatusCompleted) {
					eventType = "complete"
				} else if update.Status == string(domain.StatusFailed) {
					eventType = "error"
				}

				// Send SSE event
				data := fiber.Map{
					"status":   update.Status,
					"progress": update.Progress,
					"message":  update.Message,
				}

				if err := h.writeSSEEvent(w, eventType, data); err != nil {
					h.logger.Error("failed to write SSE event", zap.Error(err))
					return
				}

				h.logger.Debug("SSE sent event",
					zap.String("request_id", requestID.String()),
					zap.String("event_type", eventType),
					zap.Int("progress", update.Progress),
				)

				// If completed or failed, close the stream
				if eventType == "complete" || eventType == "error" {
					return
				}

			case <-ticker.C:
				// Send keep-alive comment
				if _, err := fmt.Fprintf(w, ": keepalive\n\n"); err != nil {
					return
				}
				w.Flush()
			}
		}
	})

	return nil
}

// sendSSEEvent sends an SSE event via Fiber context
func (h *EventsHandler) sendSSEEvent(c *fiber.Ctx, event string, data fiber.Map) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("event: %s\ndata: %s\n\n", event, jsonData)
	c.Write([]byte(message))
	return nil
}

// writeSSEEvent writes an SSE event to a writer
func (h *EventsHandler) writeSSEEvent(w *bufio.Writer, event string, data fiber.Map) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, jsonData); err != nil {
		return err
	}

	w.Flush()
	return nil
}
