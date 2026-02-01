package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/P4ST4S/manga-translator/backend-api/internal/domain"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type requestRepository struct {
	db *pgxpool.Pool
}

// NewRequestRepository creates a new PostgreSQL request repository
func NewRequestRepository(db *pgxpool.Pool) ports.RequestRepository {
	return &requestRepository{db: db}
}

func (r *requestRepository) Create(ctx context.Context, request *domain.Request) error {
	query := `
		INSERT INTO requests (id, filename, file_type, status, progress, page_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		request.ID,
		request.Filename,
		request.FileType,
		request.Status,
		request.Progress,
		request.PageCount,
		request.CreatedAt,
		request.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return nil
}

func (r *requestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Request, error) {
	query := `
		SELECT id, filename, file_type, status, progress, page_count,
		       thumbnail_path, error_message, created_at, updated_at, completed_at
		FROM requests
		WHERE id = $1
	`

	var request domain.Request
	err := r.db.QueryRow(ctx, query, id).Scan(
		&request.ID,
		&request.Filename,
		&request.FileType,
		&request.Status,
		&request.Progress,
		&request.PageCount,
		&request.ThumbnailPath,
		&request.ErrorMessage,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get request: %w", err)
	}

	return &request, nil
}

func (r *requestRepository) List(ctx context.Context, filter ports.RequestFilter) ([]*domain.Request, int, error) {
	// Build query with filters
	query := `
		SELECT id, filename, file_type, status, progress, page_count,
		       thumbnail_path, error_message, created_at, updated_at, completed_at
		FROM requests
	`

	countQuery := `SELECT COUNT(*) FROM requests`
	args := []interface{}{}
	argIndex := 1

	if filter.Status != nil {
		query += fmt.Sprintf(" WHERE status = $%d", argIndex)
		countQuery += fmt.Sprintf(" WHERE status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
		argIndex++
	}

	// Get total count
	var total int
	countArgs := args[:0]
	if filter.Status != nil {
		countArgs = args[:1]
	}
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count requests: %w", err)
	}

	// Get requests
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list requests: %w", err)
	}
	defer rows.Close()

	requests := []*domain.Request{}
	for rows.Next() {
		var req domain.Request
		err := rows.Scan(
			&req.ID,
			&req.Filename,
			&req.FileType,
			&req.Status,
			&req.Progress,
			&req.PageCount,
			&req.ThumbnailPath,
			&req.ErrorMessage,
			&req.CreatedAt,
			&req.UpdatedAt,
			&req.CompletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan request: %w", err)
		}
		requests = append(requests, &req)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating requests: %w", err)
	}

	return requests, total, nil
}

func (r *requestRepository) Update(ctx context.Context, request *domain.Request) error {
	query := `
		UPDATE requests
		SET filename = $1, file_type = $2, status = $3, progress = $4,
		    page_count = $5, thumbnail_path = $6, error_message = $7,
		    updated_at = $8, completed_at = $9
		WHERE id = $10
	`

	result, err := r.db.Exec(ctx, query,
		request.Filename,
		request.FileType,
		request.Status,
		request.Progress,
		request.PageCount,
		request.ThumbnailPath,
		request.ErrorMessage,
		request.UpdatedAt,
		request.CompletedAt,
		request.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *requestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus, progress int) error {
	query := `
		UPDATE requests
		SET status = $1, progress = $2, updated_at = NOW()
		WHERE id = $3
	`

	result, err := r.db.Exec(ctx, query, status, progress, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}
