# Engineering Specifications — Penster

## Architecture Overview

Penster follows **Clean Architecture** with **CQRS (Command Query Responsibility Segregation)**.

## Request Flow

```mermaid
flowchart LR
    A["HTTP Request"] --> B["Router\nrouter.go"]
    B --> C["Middleware Chain\nOTEL Tracing → Logging → Recovery"]
    C --> D["Handler\nhandler/*.go"]
    D --> E["Service\nservice/*.go"]
    E --> F["Command / Query\ncommand/*.go, query/*.go"]
    F --> G["Repository\nrepository/*.go"]
    G --> H["Database\ninfrastructure/database/query/*.sql.go"]
    H --> I["PostgreSQL"]
```

## CQRS Pattern

```mermaid
flowchart TD
    subgraph Commands["Write Path"]
        C1["Create"]
        C2["Update"]
        C3["Delete"]
        C4["UpdateBalance"]
    end
    subgraph Queries["Read Path"]
        Q1["GetByID"]
        Q2["List"]
        Q3["Summary"]
        Q4["ByAccount"]
        Q5["ByCategory"]
        Q6["Trends"]
    end
    subgraph Services["Service"]
        Svc["Orchestrates\nbusiness logic"]
    end
    C1 & C2 & C3 & C4 --> Svc
    Q1 & Q2 & Q3 & Q4 & Q5 & Q6 --> Svc
```

## Database Schema

```mermaid
erDiagram
    accounts {
        int32 id PK
        uuid sub_id UK
        string name
        string type
        numeric balance
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }
    categories {
        int32 id PK
        uuid sub_id UK
        string name
        string type
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }
    transactions {
        int32 id PK
        uuid sub_id UK
        int32 account_id FK
        int32 transfer_account_id FK
        int32 category_id FK
        string transaction_type
        string title
        numeric amount
        string currency
        numeric currency_rate
        numeric enhanced_amount
        string notes
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }
    drafts {
        int32 id PK
        uuid sub_id UK
        int32 account_id FK
        int32 transfer_account_id FK
        int32 category_id FK
        string transaction_type
        string title
        numeric amount
        string currency
        numeric currency_rate
        numeric enhanced_amount
        string notes
        string source
        string status
        timestamp confirmed_at
        timestamp rejected_at
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }
    rate_currencies {
        string from_currency PK
        string to_currency PK
        numeric rate
        date rate_date
    }

    accounts ||--o{ transactions : "source account"
    accounts ||--o{ transactions : "transfer account"
    accounts ||--o{ drafts : "source account"
    accounts ||--o{ drafts : "transfer account"
    categories ||--o{ transactions : ""
    categories ||--o{ drafts : ""
```

## Key Design Decisions

### Soft Deletes
All major entities have a `deleted_at` timestamp. Deletes set this field rather than removing rows, preserving audit history.

### Balance Update Pattern
When a transaction affects an account balance, the update follows a **"reverse first, then apply"** pattern:
1. Reverse the previous amount from the account
2. Apply the new amount

This prevents race conditions when multiple transactions update the same account concurrently.

### Parallel Validation with `syncerr.Group`
Service operations that validate multiple entities in parallel use `syncerr.Group` to collect all errors before failing, providing complete feedback rather than failing on the first error.

### UUID Sub-IDs
External-facing IDs use UUID (`sub_id`) for safe URL exposure. Internal joins use auto-increment `int32` for efficiency.

## Middleware Chain

```mermaid
sequenceDiagram
    participant Client
    participant Tracing as OTEL Tracing
    participant Logging as Logging
    participant Recovery as Recovery
    participant Handler

    Client->>Tracing: HTTP Request
    Tracing->>Logging: Wrapped Handler
    Logging->>Recovery: Wrapped Handler
    Recovery->>Handler: Wrapped Handler
    Handler-->>Recovery: Response
    Recovery-->>Logging: Response + panic catch
    Logging-->>Tracing: Response + duration
    Tracing-->>Client: Response
```

## Scheduler

The scheduler runs on a ticker-based engine (1-second interval). Jobs are dispatched when their next run time is reached.

**Current jobs:**
- `RateCurrencyJob` — fetches ECB FX rates hourly

```mermaid
flowchart LR
    A["Scheduler Engine\n(ticker every 1s)"] --> B{"Job due?"}
    B -->|Yes| C["RateCurrencyJob"]
    B -->|No| D["Skip"]
    C --> E["Fetch ECB rates\nUpsert rate_currencies"]
```

## Observability

OpenTelemetry spans are created at three layers:
- **HTTP layer** — `otelhttp.NewHandler` wraps all requests
- **Service layer** — `observability.StartServiceSpan()`
- **Repository layer** — `observability.StartRepoSpan()`

Traces are exported via OTLP to Jaeger.

## Layer Responsibilities

| Layer | Responsibility |
|---|---|
| `domain/entities` | Core business objects, no dependencies |
| `domain/repository` | Interface definitions (no implementation) |
| `application/command` | Write operations (Create, Update, Delete) |
| `application/query` | Read operations (GetByID, List, reports) |
| `application/service` | Business logic orchestration |
| `interface/handler` | HTTP request/response handling |
| `interface/dto` | Request validation, response transformation |
| `interface/router` | Route registration |
| `infrastructure/database` | sqlc-generated SQL queries |
| `infrastructure/postgres` | Connection pool, migrations |
| `pkg/observability` | OTEL setup, span helpers |