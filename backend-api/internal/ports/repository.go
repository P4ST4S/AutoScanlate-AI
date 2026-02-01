package ports

import (
	"context"

	"github.com/P4ST4S/manga-translator/backend-api/internal/domain"
	"github.com/google/uuid"
)

// RequestRepository defines the interface for request data persistence
type RequestRepository interface {
	// Create creates a new request
	Create(ctx context.Context, request *domain.Request) error

	// GetByID retrieves a request by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Request, error)

	// List retrieves requests with optional filtering
	List(ctx context.Context, filter RequestFilter) ([]*domain.Request, int, error)

	// Update updates an existing request
	Update(ctx context.Context, request *domain.Request) error

	// UpdateStatus updates request status and progress
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus, progress int) error
}

// ResultRepository defines the interface for result data persistence
type ResultRepository interface {
	// Create creates a new result
	Create(ctx context.Context, result *domain.Result) error

	// CreateBatch creates multiple results
	CreateBatch(ctx context.Context, results []*domain.Result) error

	// GetByRequestID retrieves all results for a request
	GetByRequestID(ctx context.Context, requestID uuid.UUID) ([]*domain.Result, error)

	// DeleteByRequestID deletes all results for a request
	DeleteByRequestID(ctx context.Context, requestID uuid.UUID) error
}

// RequestFilter represents filtering options for listing requests
type RequestFilter struct {
	Status *domain.RequestStatus
	Limit  int
	Offset int
}
