# Product Requirements Document — Penster

> Personal Finance Transaction Management API

## Overview

**Penster** is a Go-based personal finance API built with Clean Architecture and CQRS. It manages accounts, categories, transactions, and drafts with full OpenTelemetry observability.

## Core Entities

### Account
- Holds a balance in a specific currency
- Types: `expense`, `income`, `transfer`
- Supports soft delete
- External ID uses UUID (`sub_id`), internal ID uses `int32`

### Category
- Classifies transactions
- Types: `expense`, `income`, `transfer`
- Supports soft delete

### Transaction
- Represents a confirmed financial movement
- Links to one account, optionally a transfer account, and a category
- Types: `expense`, `income`, `transfer`
- Stores `amount`, `currency`, `currency_rate`, and `enhanced_amount` (converted to base currency)
- Supports soft delete

### Draft
- A **pending transaction** awaiting confirmation or rejection
- Draft lifecycle:
  ```
  [Created] --> [Confirm] --> [Creates Transaction]
                or
  [Created] --> [Reject] --> [Marked Rejected]
  ```
- Sources: `manual` (API) or `ingestion` (future import)
- Statuses: `pending`, `confirmed`, `rejected`

### RateCurrency
- Stores ECB foreign exchange rates
- Fetches daily from ECB API
- Used to convert transaction amounts to base currency

## Data Model Decisions

| Decision | Rationale |
|---|---|
| UUID `sub_id` for external IDs | Safe to expose in URLs, no enumeration |
| Internal `int32` ID for joins | Efficient, auto-increment |
| Soft deletes (`deleted_at`) | Audit trail, no accidental data loss |
| Balance "reverse then apply" | Prevents race conditions on concurrent updates |
| `syncerr.Group` for parallel validation | Collect all errors before failing |

## Reports

- **Summary** — total balance, total expenses, total income, total transfers
- **By Account** — breakdown per account
- **By Category** — breakdown per category
- **Trends** — data points over time grouped by date and type

## Technical Constraints

- **Language**: Go 1.25
- **Database**: PostgreSQL 16
- **Code Generation**: sqlc (SQL-first queries)
- **Observability**: OpenTelemetry (OTLP to Jaeger)
- **API Documentation**: Swagger at `/swagger/`
- **Authentication**: None (public API — auth TBD)

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `APP_ENV` | `local` | Environment name |
| `APP_PORT` | `8080` | HTTP server port |
| `APP_VERSION` | `1.0.0` | App version |
| `BASE_CURRENCY` | `IDR` | Base currency for reports |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `penster` | Database user |
| `DB_PASSWORD` | `placeholder` | Database password |
| `DB_NAME` | `penster` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `DB_MAX_CONNS` | `10` | Max connections |
| `DB_MIN_CONNS` | `2` | Min connections |
| `AUTO_MIGRATE` | `true` | Run migrations on startup |
| `OTEL_ENABLED` | `true` | Enable OTEL tracing |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `jaeger:4317` | OTEL collector endpoint |
| `ECB_URL` | ECB daily XML | Exchange rate source |