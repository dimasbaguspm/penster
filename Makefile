.PHONY: init dev-backend dev-reset

init:
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod tidy
	go mod download

dev-backend:
	docker compose -f infra/docker-compose.local.yml up --build

dev-reset:
	docker compose -f infra/docker-compose.local.yml down -v
	docker compose -f infra/docker-compose.local.yml up --build
