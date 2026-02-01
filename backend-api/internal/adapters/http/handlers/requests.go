package handlers

import (
	"errors"

	"github.com/P4ST4S/manga-translator/backend-api/internal/domain"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RequestsHandler struct {
	requestRepo ports.RequestRepository
	logger      *zap.Logger
}

func NewRequestsHandler(requestRepo ports.RequestRepository, logger *zap.Logger) *RequestsHandler {
	return &RequestsHandler{
		requestRepo: requestRepo,
		logger:      logger,
	}
}

// List handles GET /api/requests
func (h *RequestsHandler) List(c *fiber.Ctx) error {
	// Parse query parameters
	var filter ports.RequestFilter
	filter.Limit = c.QueryInt("limit", 20)
	filter.Offset = c.QueryInt("offset", 0)

	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.RequestStatus(statusStr)
		filter.Status = &status
	}

	// Get requests from repository
	requests, total, err := h.requestRepo.List(c.Context(), filter)
	if err != nil {
		h.logger.Error("failed to list requests", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve requests",
		})
	}

	return c.JSON(fiber.Map{
		"requests": requests,
		"total":    total,
		"limit":    filter.Limit,
		"offset":   filter.Offset,
	})
}

// GetByID handles GET /api/requests/:id
func (h *RequestsHandler) GetByID(c *fiber.Ctx) error {
	// Parse ID
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request ID",
		})
	}

	// Get request from repository
	request, err := h.requestRepo.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "request not found",
			})
		}
		h.logger.Error("failed to get request", zap.Error(err), zap.String("id", idStr))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve request",
		})
	}

	return c.JSON(request)
}
