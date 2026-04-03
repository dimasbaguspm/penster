.PHONY: init dev-backend dev-reset swag

init:
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod tidy
	go mod download

dev-backend:
	docker compose -f infra/docker-compose.local.yml up --build

dev-reset:
	docker compose -f infra/docker-compose.local.yml down -v

swag:
	swag init -g ./cmd/server/main.go -o ./docs --packageName docs --quiet --generatedTime
