package handlers

import (
	"errors"

	"github.com/P4ST4S/manga-translator/backend-api/internal/domain"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ResultsHandler struct {
	requestRepo ports.RequestRepository
	resultRepo  ports.ResultRepository
	logger      *zap.Logger
}

func NewResultsHandler(
	requestRepo ports.RequestRepository,
	resultRepo ports.ResultRepository,
	logger *zap.Logger,
) *ResultsHandler {
	return &ResultsHandler{
		requestRepo: requestRepo,
		resultRepo:  resultRepo,
		logger:      logger,
	}
}

// GetByRequestID handles GET /api/results/:id
func (h *ResultsHandler) GetByRequestID(c *fiber.Ctx) error {
	// Parse ID
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request ID",
		})
	}

	// Check if request exists
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

	// Check if request is completed
	if !request.IsCompleted() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "request not yet completed",
		})
	}

	// Get results
	results, err := h.resultRepo.GetByRequestID(c.Context(), id)
	if err != nil {
		h.logger.Error("failed to get results", zap.Error(err), zap.String("requestId", idStr))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve results",
		})
	}

	return c.JSON(fiber.Map{
		"requestId": id,
		"pages":     results,
	})
}
