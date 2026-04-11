# Penster

Personal Finance Transaction Management API

[![Go Version](https://img.shields.io/badge/Go-1.25-blue)](https://golang.org/)
[![Docker Image](https://img.shields.io/badge/Docker-dimasbaguspm%2Fpenster:latest-blue?logo=docker)](https://github.com/dimasbaguspm/penster/pkgs/container/penster)

## Quick Start

```bash
# Clone
git clone https://github.com/dimasbaguspm/penster.git
cd penster

# Start (postgres + jaeger + penster)
docker compose -f infra/docker-compose.local.yml up -d

# Check health
curl http://localhost:8080/health

# Open Swagger UI
open http://localhost:8080/swagger/
```


## Documentation

| Doc | Description |
|---|---|
| [PRD](docs/PRD.md) | Product requirements, entity definitions, technical constraints |
| [Engineering Specs](docs/engineering-specs.md) | Architecture diagrams, CQRS flow, database schema, design decisions |
| [Endpoints](docs/endpoints.md) | All 24 API routes with request/response shapes |
| [Architecture](docs/architecture.md) | Layer responsibilities, DI flow, observability, patterns |
| [Installation](docs/installation.md) | Docker setup, environment variables, database migrations |

