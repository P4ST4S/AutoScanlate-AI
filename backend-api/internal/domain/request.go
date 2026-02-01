package domain

import (
	"time"

	"github.com/google/uuid"
)

// RequestStatus represents the status of a translation request
type RequestStatus string

const (
	StatusQueued     RequestStatus = "queued"
	StatusProcessing RequestStatus = "processing"
	StatusCompleted  RequestStatus = "completed"
	StatusFailed     RequestStatus = "failed"
)

// FileType represents the type of uploaded file
type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeZip   FileType = "zip"
)

// Request represents a translation request
type Request struct {
	ID            uuid.UUID     `json:"id"`
	Filename      string        `json:"filename"`
	FileType      FileType      `json:"fileType"`
	Status        RequestStatus `json:"status"`
	Progress      int           `json:"progress"`
	PageCount     int           `json:"pageCount"`
	ThumbnailPath *string       `json:"thumbnail,omitempty"`
	ErrorMessage  *string       `json:"errorMessage,omitempty"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	CompletedAt   *time.Time    `json:"completedAt,omitempty"`
}

// NewRequest creates a new translation request
func NewRequest(filename string, fileType FileType) *Request {
	now := time.Now()
	return &Request{
		ID:        uuid.New(),
		Filename:  filename,
		FileType:  fileType,
		Status:    StatusQueued,
		Progress:  0,
		PageCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateStatus updates the request status and progress
func (r *Request) UpdateStatus(status RequestStatus, progress int) {
	r.Status = status
	r.Progress = progress
	r.UpdatedAt = time.Now()

	if status == StatusCompleted || status == StatusFailed {
		now := time.Now()
		r.CompletedAt = &now
	}
}

// SetError sets the error message and marks the request as failed
func (r *Request) SetError(message string) {
	r.Status = StatusFailed
	r.ErrorMessage = &message
	r.UpdatedAt = time.Now()
	now := time.Now()
	r.CompletedAt = &now
}

// IsCompleted returns true if the request is completed or failed
func (r *Request) IsCompleted() bool {
	return r.Status == StatusCompleted || r.Status == StatusFailed
}
