.PHONY: init vet build test swag sql dev-backend dev-reset

init:
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod tidy
	go mod download

vet:
	go vet ./...

build:
	go build -o bin/penster ./cmd/server

test:
	go clean -testcache
	go build -o bin/penster_test ./cmd/server
	go test -v -count=1 ./tests/...

swag:
	swag init -g ./cmd/server/main.go -o ./docs --packageName docs --quiet

sql:
	sqlc generate

dev-backend:
	docker compose -f infra/docker-compose.local.yml up --build

dev-reset:
	docker compose -f infra/docker-compose.local.yml down -v
