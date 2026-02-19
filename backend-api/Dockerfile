# Build stage
# Build context is the project root (set in docker-compose.yml)
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend-api/go.mod backend-api/go.sum ./
RUN go mod download

# Copy source code
COPY backend-api/ .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/api .

# Copy migrations
COPY backend-api/migrations ./migrations

# Create storage directories
RUN mkdir -p storage/uploads storage/originals storage/translated storage/temp

# Expose port
EXPOSE 8080

# Run the application
CMD ["./api"]
