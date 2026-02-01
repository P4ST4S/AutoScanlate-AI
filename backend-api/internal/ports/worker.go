package ports

import (
	"context"
)

// ProgressCallback is called when translation progress is updated
type ProgressCallback func(progress int, message string)

// TranslationOutput represents the output of a translation job
type TranslationOutput struct {
	OutputPath string
	Pages      []PageOutput
}

// PageOutput represents a single translated page
type PageOutput struct {
	PageNumber     int
	OriginalPath   string
	TranslatedPath string
}

// WorkerExecutor defines the interface for executing translation jobs
type WorkerExecutor interface {
	// Translate executes a translation job
	Translate(ctx context.Context, inputPath string, onProgress ProgressCallback) (*TranslationOutput, error)
}
