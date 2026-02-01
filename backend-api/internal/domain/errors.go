package domain

import "errors"

var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrUploadFailed is returned when file upload fails
	ErrUploadFailed = errors.New("file upload failed")

	// ErrWorkerFailed is returned when the translation worker fails
	ErrWorkerFailed = errors.New("translation worker failed")

	// ErrDatabaseError is returned when a database operation fails
	ErrDatabaseError = errors.New("database operation failed")

	// ErrQueueError is returned when a queue operation fails
	ErrQueueError = errors.New("queue operation failed")

	// ErrInvalidFileType is returned when an unsupported file type is uploaded
	ErrInvalidFileType = errors.New("invalid file type")

	// ErrFileTooLarge is returned when uploaded file exceeds size limit
	ErrFileTooLarge = errors.New("file too large")

	// ErrTooManyFiles is returned when too many files are uploaded
	ErrTooManyFiles = errors.New("too many files")
)

// AppError represents an application error with additional context
type AppError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Code + ": " + e.Message + " - " + e.Err.Error()
	}
	return e.Code + ": " + e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
