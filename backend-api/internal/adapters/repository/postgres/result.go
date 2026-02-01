package postgres

import (
	"context"
	"fmt"

	"github.com/P4ST4S/manga-translator/backend-api/internal/domain"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type resultRepository struct {
	db *pgxpool.Pool
}

// NewResultRepository creates a new PostgreSQL result repository
func NewResultRepository(db *pgxpool.Pool) ports.ResultRepository {
	return &resultRepository{db: db}
}

func (r *resultRepository) Create(ctx context.Context, result *domain.Result) error {
	query := `
		INSERT INTO results (id, request_id, page_number, original_path, translated_path, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		result.ID,
		result.RequestID,
		result.PageNumber,
		result.OriginalPath,
		result.TranslatedPath,
		result.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create result: %w", err)
	}

	return nil
}

func (r *resultRepository) CreateBatch(ctx context.Context, results []*domain.Result) error {
	if len(results) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO results (id, request_id, page_number, original_path, translated_path, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, result := range results {
		_, err := tx.Exec(ctx, query,
			result.ID,
			result.RequestID,
			result.PageNumber,
			result.OriginalPath,
			result.TranslatedPath,
			result.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert result: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *resultRepository) GetByRequestID(ctx context.Context, requestID uuid.UUID) ([]*domain.Result, error) {
	query := `
		SELECT id, request_id, page_number, original_path, translated_path, created_at
		FROM results
		WHERE request_id = $1
		ORDER BY page_number ASC
	`

	rows, err := r.db.Query(ctx, query, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get results: %w", err)
	}
	defer rows.Close()

	results := []*domain.Result{}
	for rows.Next() {
		var result domain.Result
		err := rows.Scan(
			&result.ID,
			&result.RequestID,
			&result.PageNumber,
			&result.OriginalPath,
			&result.TranslatedPath,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan result: %w", err)
		}
		results = append(results, &result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return results, nil
}

func (r *resultRepository) DeleteByRequestID(ctx context.Context, requestID uuid.UUID) error {
	query := `DELETE FROM results WHERE request_id = $1`

	_, err := r.db.Exec(ctx, query, requestID)
	if err != nil {
		return fmt.Errorf("failed to delete results: %w", err)
	}

	return nil
}
