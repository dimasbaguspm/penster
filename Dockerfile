FROM golang:1.25-alpine AS dev
RUN apk add --no-cache ca-certificates tzdata git make && \
    go install github.com/air-verse/air@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /app
CMD ["/go/bin/air"]