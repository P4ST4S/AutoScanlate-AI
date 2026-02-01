# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates python3 py3-pip tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/api .

# Copy migrations
COPY migrations ./migrations

# Create storage directories
RUN mkdir -p storage/uploads storage/originals storage/translated storage/temp

# Expose port
EXPOSE 8080

# Run the application
CMD ["./api"]
