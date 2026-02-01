package ports

import (
	"context"
	"io"
)

// Storage defines the interface for file storage operations
type Storage interface {
	// Save saves a file and returns the storage path
	Save(ctx context.Context, path string, data io.Reader) error

	// Get retrieves a file
	Get(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete deletes a file
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists
	Exists(ctx context.Context, path string) (bool, error)
}
