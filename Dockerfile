# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Generate SQLC code
RUN which sqlc || go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN sqlc generate

# Final stage
FROM alpine:3.19 AS runner

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary and generated code from builder
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/infrastructure/database/query ./query
COPY --from=builder /app/sqlc.yaml .

# Copy config
COPY config ./config

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
