# Installation — Penster

## Prerequisites

- Docker & Docker Compose
- PostgreSQL 16 (if not using Docker)

---

## Quick Start (Docker Compose)

```bash
# Clone the repository
git clone https://github.com/dimasbaguspm/penster.git
cd penster

# Start the stack (postgres + jaeger + penster)
docker compose -f infra/docker-compose.local.yml up -d

# Check health
curl http://localhost:8080/health
# {"status":"ok","timestamp":"2026-04-11T10:00:00Z","version":"1.0.0"}

# Access Swagger UI
open http://localhost:8080/swagger/
```

---

## Pre-built Docker Image

The official image is hosted at GitHub Container Registry:

```
docker.io/dimasbaguspm/penster:latest
```

```bash
# Pull the image
docker pull docker.io/dimasbaguspm/penster:latest

# Run with environment variables
docker run -d \
  --name penster \
  -p 8080:8080 \
  -e DB_HOST=your-postgres-host \
  -e DB_PORT=5432 \
  -e DB_USER=penster \
  -e DB_PASSWORD=your-password \
  -e DB_NAME=penster \
  docker.io/dimasbaguspm/penster:latest
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `APP_ENV` | `local` | Environment (`local`, `production`, etc.) |
| `APP_PORT` | `8080` | HTTP server port |
| `APP_VERSION` | `1.0.0` | App version string |
| `BASE_CURRENCY` | `IDR` | Base currency for report conversions |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `penster` | Database user |
| `DB_PASSWORD` | `placeholder` | Database password |
| `DB_NAME` | `penster` | Database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `DB_MAX_CONNS` | `10` | Max connections in pool |
| `DB_MIN_CONNS` | `2` | Min connections in pool |
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker addresses |
| `AUTO_MIGRATE` | `true` | Run SQL migrations on startup |
| `MIGRATE_PATH` | `migrations` | Path to migration files |
| `OTEL_ENABLED` | `true` | Enable OpenTelemetry tracing |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `jaeger:4317` | OTEL collector gRPC address |
| `ECB_URL` | ECB daily XML URL | Exchange rate source |

Create a `.env` file from the example:

```bash
cp .env.example .env
# Edit .env with your values
```

---

## Database Setup

### Automatic (Default)

Set `AUTO_MIGRATE=true` (default). Migrations run automatically on startup.

### Manual

If `AUTO_MIGRATE=false`, run migrations manually:

```bash
# Using golang-migrate
migrate -path ./migrations -database "postgres://penster:password@localhost:5432/penster?sslmode=disable" up

# Or via Docker
docker run -it --rm \
  -e DATABASE_URL="postgres://penster:password@localhost:5432/penster?sslmode=disable" \
  migrate up
```

### Migrations

| ID | Name | Description |
|---|---|---|
| 000001 | init_schema | accounts, categories tables |
| 000002 | add_rate_currencies | currency rate table |
| 000003 | add_transactions | transactions table |
| 000004 | add_drafts | drafts (pending transactions) |
| 000005 | add_report_indexes | performance indexes |

---

## Docker Compose (Local Development)

`infra/docker-compose.local.yml`:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: penster
      POSTGRES_PASSWORD: placeholder
      POSTGRES_DB: penster
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  jaeger:
    image: jaeger:2.17.0
    ports:
      - "16686:16686"  # UI
      - "16685:16685"  # OTLP gRPC
      - "4317:4317"    # OTLP HTTP

  penster:
    build: .
    ports:
      - "8080:8080"
    environment:
      APP_ENV: local
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: penster
      DB_PASSWORD: placeholder
      DB_NAME: penster
      AUTO_MIGRATE: "true"
      OTEL_EXPORTER_OTLP_ENDPOINT: jaeger:4317
    depends_on:
      - postgres
      - jaeger

volumes:
  postgres_data:
```

```bash
# Start
docker compose -f infra/docker-compose.local.yml up -d

# Stop
docker compose -f infra/docker-compose.local.yml down
```

---

## Build from Source (Development)

Requires: Go 1.25+

```bash
# Install dependencies
go mod download

# Run with live reload (air)
air

# Or run directly
go run cmd/server/main.go
```

### Development Tools

| Tool | Purpose |
|---|---|
| `air` | Live reload (see `air.toml`) |
| `sqlc` | SQL-first code generation (see `sqlc.yaml`) |

```bash
# Regenerate database queries after SQL changes
sqlc generate
```

---

## Health Checks

Once running:

```bash
# Liveness
curl http://localhost:8080/health

# Readiness
curl http://localhost:8080/ready
```

---

## Ports Summary

| Port | Service |
|---|---|
| `8080` | Penster API |
| `8080/swagger/` | Swagger UI |
| `5432` | PostgreSQL |
| `16686` | Jaeger UI |
| `16685` | Jaeger OTLP gRPC |
| `4317` | OTLP HTTP |

---

## Upgrading

1. Pull the new image:
   ```bash
   docker pull docker.io/dimasbaguspm/penster:latest
   ```
2. Restart the container:
   ```bash
   docker compose -f infra/docker-compose.local.yml up -d
   ```
3. Migrations run automatically (`AUTO_MIGRATE=true`).