package handlers

import (
	"path/filepath"
	"strings"

	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type FilesHandler struct {
	storagePath string
	logger      *zap.Logger
}

func NewFilesHandler(cfg *config.Config, logger *zap.Logger) *FilesHandler {
	return &FilesHandler{
		storagePath: cfg.Storage.Path,
		logger:      logger,
	}
}

// ServeFile handles GET /api/files/:requestId/:type/*
func (h *FilesHandler) ServeFile(c *fiber.Ctx) error {
	requestID := c.Params("requestId")
	fileType := c.Params("type")
	filePath := c.Params("*") // Capture the remaining path (can include subdirectories)

	// Validate type
	validTypes := map[string]bool{
		"uploads":    true,
		"originals":  true,
		"translated": true,
	}

	if !validTypes[fileType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid file type",
		})
	}

	// Sanitize file path to prevent path traversal
	if strings.Contains(filePath, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid file path",
		})
	}

	// Clean the path to normalize slashes
	filePath = filepath.Clean(filePath)

	// Build full file path
	fullPath := filepath.Join(h.storagePath, fileType, requestID, filePath)

	// Serve file
	return c.SendFile(fullPath)
}
