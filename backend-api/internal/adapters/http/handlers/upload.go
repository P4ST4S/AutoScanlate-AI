package handlers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/P4ST4S/manga-translator/backend-api/internal/domain"
	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/config"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UploadHandler struct {
	requestRepo ports.RequestRepository
	queueClient ports.QueueClient
	cfg         *config.Config
	logger      *zap.Logger
}

func NewUploadHandler(
	requestRepo ports.RequestRepository,
	queueClient ports.QueueClient,
	cfg *config.Config,
	logger *zap.Logger,
) *UploadHandler {
	return &UploadHandler{
		requestRepo: requestRepo,
		queueClient: queueClient,
		cfg:         cfg,
		logger:      logger,
	}
}

// Upload handles POST /api/translate
func (h *UploadHandler) Upload(c *fiber.Ctx) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse form data",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no files provided",
		})
	}

	if len(files) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "too many files (max 10)",
		})
	}

	// Process first file (simplified for now)
	file := files[0]

	// Validate file size
	if file.Size > h.cfg.Storage.MaxUploadSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file too large",
		})
	}

	// Validate file type
	filename := file.Filename
	fileType := h.getFileType(filename)
	if fileType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unsupported file type",
		})
	}

	// Create request
	request := domain.NewRequest(filename, domain.FileType(fileType))

	// Create upload directory
	uploadDir := filepath.Join(h.cfg.Storage.Path, "uploads", request.ID.String())
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		h.logger.Error("failed to create upload directory", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save file",
		})
	}

	// Save file
	filePath := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		h.logger.Error("failed to save file", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save file",
		})
	}

	// Save request to database
	if err := h.requestRepo.Create(c.Context(), request); err != nil {
		h.logger.Error("failed to create request", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create request",
		})
	}

	// Enqueue translation job
	if h.queueClient != nil {
		if err := h.queueClient.EnqueueTranslation(c.Context(), request.ID, filePath, fileType); err != nil {
			h.logger.Error("failed to enqueue job", zap.Error(err))
			// Continue anyway - job can be retried
		}
	}

	h.logger.Info("file uploaded",
		zap.String("requestId", request.ID.String()),
		zap.String("filename", filename),
		zap.Int64("size", file.Size),
	)

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (h *UploadHandler) getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if imageExts[ext] {
		return "image"
	}

	if ext == ".zip" {
		return "zip"
	}

	return ""
}
