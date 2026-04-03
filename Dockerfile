FROM golang:1.25-alpine AS dev
RUN apk add ca-certificates tzdata git make && \
    go install github.com/air-verse/air@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /app

FROM golang:1.25-alpine AS builder
RUN apk add --no-cache ca-certificates tzdata git make
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /app/bin/server ./cmd/server

COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["/go/bin/air"]
