package domain

import (
	"time"

	"github.com/google/uuid"
)

// Result represents a translated page result
type Result struct {
	ID             uuid.UUID `json:"id"`
	RequestID      uuid.UUID `json:"requestId"`
	PageNumber     int       `json:"pageNumber"`
	OriginalPath   string    `json:"original"`
	TranslatedPath string    `json:"translated"`
	CreatedAt      time.Time `json:"createdAt"`
}

// NewResult creates a new result for a translated page
func NewResult(requestID uuid.UUID, pageNumber int, originalPath, translatedPath string) *Result {
	return &Result{
		ID:             uuid.New(),
		RequestID:      requestID,
		PageNumber:     pageNumber,
		OriginalPath:   originalPath,
		TranslatedPath: translatedPath,
		CreatedAt:      time.Now(),
	}
}

// ResultList represents a collection of results for a request
type ResultList struct {
	RequestID uuid.UUID `json:"requestId"`
	Pages     []Result  `json:"pages"`
}
