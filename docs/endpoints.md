# API Endpoints — Penster

Base URL: `http://localhost:8080`
Swagger UI: `http://localhost:8080/swagger/`

**Note:** No authentication — all endpoints are publicly accessible.

---

## Health

| Method | Path | Handler | Description |
|---|---|---|---|
| `GET` | `/health` | `healthHandler.Health` | Returns `{status, timestamp, version}` |
| `GET` | `/ready` | `healthHandler.Ready` | Readiness probe |

---

## Accounts

| Method | Path | Handler | Description |
|---|---|---|---|
| `GET` | `/accounts` | `accountHandler.List` | List all accounts (no pagination) |
| `POST` | `/accounts` | `accountHandler.Create` | Create an account |
| `GET` | `/accounts/{id}` | `accountHandler.Get` | Get account by `sub_id` |
| `PUT` | `/accounts/{id}` | `accountHandler.Update` | Update account |
| `DELETE` | `/accounts/{id}` | `accountHandler.Delete` | Soft delete account |

### Account Models

**Account:**
```json
{
  "id": 1,
  "sub_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Bank Account",
  "type": "expense",
  "balance": 1500000,
  "deleted_at": null,
  "created_at": "2026-04-11T10:00:00Z",
  "updated_at": "2026-04-11T10:00:00Z"
}
```

**CreateAccountRequest:**
```json
{
  "name": "Bank Account",
  "type": "expense",
  "balance": 1500000
}
```

**UpdateAccountRequest:**
```json
{
  "name": "Savings Account",
  "type": "income",
  "balance": 2000000
}
```

---

## Categories

| Method | Path | Handler | Description |
|---|---|---|---|
| `GET` | `/categories` | `categoryHandler.List` | List all categories |
| `POST` | `/categories` | `categoryHandler.Create` | Create a category |
| `GET` | `/categories/{id}` | `categoryHandler.Get` | Get category by `sub_id` |
| `PUT` | `/categories/{id}` | `categoryHandler.Update` | Update category |
| `DELETE` | `/categories/{id}` | `categoryHandler.Delete` | Soft delete category |

### Category Models

**Category:**
```json
{
  "id": 1,
  "sub_id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "Groceries",
  "type": "expense",
  "deleted_at": null,
  "created_at": "2026-04-11T10:00:00Z",
  "updated_at": "2026-04-11T10:00:00Z"
}
```

**CreateCategoryRequest:**
```json
{
  "name": "Groceries",
  "type": "expense"
}
```

**UpdateCategoryRequest:**
```json
{
  "name": "Food & Groceries",
  "type": "expense"
}
```

`type` must be one of: `expense`, `income`, `transfer`

---

## Transactions

| Method | Path | Handler | Description |
|---|---|---|---|
| `GET` | `/transactions` | `transactionHandler.List` | List transactions |
| `POST` | `/transactions` | `transactionHandler.Create` | Create a transaction |
| `GET` | `/transactions/{id}` | `transactionHandler.Get` | Get transaction by `sub_id` |
| `PUT` | `/transactions/{id}` | `transactionHandler.Update` | Update transaction |
| `DELETE` | `/transactions/{id}` | `transactionHandler.Delete` | Soft delete transaction |

### Transaction Models

**Transaction:**
```json
{
  "id": 1,
  "sub_id": "550e8400-e29b-41d4-a716-446655440002",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "transfer_account_id": null,
  "category_id": "550e8400-e29b-41d4-a716-446655440001",
  "transaction_type": "expense",
  "title": "Weekly groceries",
  "amount": 250000,
  "currency": "IDR",
  "currency_rate": 1,
  "enhanced_amount": 250000,
  "notes": "Weekly shopping",
  "deleted_at": null,
  "created_at": "2026-04-11T10:00:00Z",
  "updated_at": "2026-04-11T10:00:00Z"
}
```

**CreateTransactionRequest:**
```json
{
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "transfer_account_id": null,
  "category_id": "550e8400-e29b-41d4-a716-446655440001",
  "transaction_type": "expense",
  "title": "Weekly groceries",
  "amount": 250000,
  "currency": "IDR",
  "notes": "Weekly shopping"
}
```

`transaction_type` must be one of: `expense`, `income`, `transfer`
For `transfer`, `transfer_account_id` is required.

---

## Drafts

| Method | Path | Handler | Description |
|---|---|---|---|
| `GET` | `/drafts` | `draftHandler.List` | List drafts |
| `POST` | `/drafts` | `draftHandler.Create` | Create a draft |
| `GET` | `/drafts/{id}` | `draftHandler.Get` | Get draft by `sub_id` |
| `PATCH` | `/drafts/{id}` | `draftHandler.Update` | Update draft |
| `POST` | `/drafts/{id}/confirm` | `draftHandler.Confirm` | Confirm draft → creates transaction |
| `POST` | `/drafts/{id}/reject` | `draftHandler.Reject` | Reject draft |
| `DELETE` | `/drafts/{id}` | `draftHandler.Delete` | Soft delete draft |

### Draft Lifecycle

```
Draft (pending) --> Confirm --> Transaction created
       |
       --> Reject --> Draft marked rejected
```

### Draft Models

**Draft:**
```json
{
  "id": 1,
  "sub_id": "550e8400-e29b-41d4-a716-446655440003",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "transfer_account_id": null,
  "category_id": "550e8400-e29b-41d4-a716-446655440001",
  "transaction_type": "expense",
  "title": "Weekly groceries",
  "amount": 250000,
  "currency": "IDR",
  "currency_rate": 1,
  "enhanced_amount": 250000,
  "notes": null,
  "source": "manual",
  "status": "pending",
  "confirmed_at": null,
  "rejected_at": null,
  "deleted_at": null,
  "created_at": "2026-04-11T10:00:00Z",
  "updated_at": "2026-04-11T10:00:00Z"
}
```

**CreateDraftRequest:**
```json
{
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "category_id": "550e8400-e29b-41d4-a716-446655440001",
  "transaction_type": "expense",
  "title": "Weekly groceries",
  "amount": 250000,
  "currency": "IDR",
  "notes": null,
  "source": "manual"
}
```

`source` must be one of: `manual`, `ingestion`
`status` is one of: `pending`, `confirmed`, `rejected`

---

## Reports

| Method | Path | Handler | Description |
|---|---|---|---|
| `GET` | `/reports/summary` | `reportHandler.Summary` | Overall financial summary |
| `GET` | `/reports/by-account` | `reportHandler.ByAccount` | Totals per account |
| `GET` | `/reports/by-category` | `reportHandler.ByCategory` | Totals per category |
| `GET` | `/reports/trends` | `reportHandler.Trends` | Data points over time |

Reports support query parameters:
- `start_date` (optional) — filter period start (format: `YYYY-MM-DD`)
- `end_date` (optional) — filter period end (format: `YYYY-MM-DD`)
- `base_currency` (optional) — override base currency

### Report Models

**ReportSummary:**
```json
{
  "data": {
    "total_balance": 5000000,
    "total_expenses": 1500000,
    "total_income": 6500000,
    "total_transfers": 0,
    "base_currency": "IDR",
    "period_start": "2026-04-01",
    "period_end": "2026-04-30"
  }
}
```

**ReportByAccount:**
```json
{
  "data": {
    "accounts": [
      {
        "account_id": "550e8400-e29b-41d4-a716-446655440000",
        "account_name": "Bank Account",
        "type": "expense",
        "total": 1500000
      }
    ],
    "period_start": "2026-04-01",
    "period_end": "2026-04-30"
  }
}
```

**ReportByCategory:**
```json
{
  "data": {
    "categories": [
      {
        "category_id": "550e8400-e29b-41d4-a716-446655440001",
        "category_name": "Groceries",
        "type": "expense",
        "total": 250000
      }
    ],
    "period_start": "2026-04-01",
    "period_end": "2026-04-30"
  }
}
```

**ReportTrends:**
```json
{
  "data": {
    "data_points": [
      {
        "date": "2026-04-10",
        "type": "expense",
        "total": 500000
      }
    ],
    "period_start": "2026-04-01",
    "period_end": "2026-04-30"
  }
}
```

---

## Response Format

**Single resource:**
```json
{
  "success": true,
  "data": { ... },
  "error": ""
}
```

**Paginated:**
```json
{
  "success": true,
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  },
  "error": ""
}
```

**Error:**
```json
{
  "success": false,
  "error": "account not found"
}
```