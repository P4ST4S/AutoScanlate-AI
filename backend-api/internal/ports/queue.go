package ports

import (
	"context"

	"github.com/google/uuid"
)

// QueueClient defines the interface for job queue operations
type QueueClient interface {
	// EnqueueTranslation enqueues a translation job
	EnqueueTranslation(ctx context.Context, requestID uuid.UUID, filePath string, fileType string) error

	// Close closes the queue client
	Close() error
}

// QueueServer defines the interface for processing jobs
type QueueServer interface {
	// Start starts the queue server
	Start() error

	// Stop stops the queue server
	Stop() error
}
