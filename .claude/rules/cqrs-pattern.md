---
name: cqrs-pattern
description: CQRS pattern for implementing new features in Penster
type: rule
paths:
  - "internal/application/**/*.go"
  - "internal/domain/**/*.go"
  - "internal/interface/**/*.go"
  - "cmd/server/**/*.go"
---

# CQRS Pattern in Penster

## When to Apply
When adding a new entity or feature, implement it across all three layers.

## Layer Structure

### 1. Domain Layer (`internal/domain/`)
- **Entities** in `entities/` — core business objects with ID, soft deletes
- **Repository interfaces** in `repository/` — define data access contracts
- **Value objects** in `valueobjects/` — parameter transformation helpers

### 2. Application Layer (`internal/application/`)

**Command** (`command/`): Write operations
```go
type TransactionCommandInterface interface {
    Create(ctx context.Context, params query.CreateTransactionParams) (*models.Transaction, error)
    Update(ctx context.Context, id string, params query.UpdateTransactionParams) (*models.Transaction, error)
    Delete(ctx context.Context, id string) (*models.Transaction, error)
}
var _ TransactionCommandInterface = (*TransactionCommand)(nil)
```

**Query** (`query/`): Read operations
```go
type TransactionQueryInterface interface {
    GetByID(ctx context.Context, id string) (*models.Transaction, error)
    List(ctx context.Context, params query.ListTransactionsParams) ([]*models.Transaction, error)
}
```

**Service** (`service/`): Business logic orchestration
- `*Service` structs hold repositories and implement business rules
- Balance update/reverse logic lives here

### 3. Interface Layer (`internal/interface/`)
- **DTOs** in `dto/` — API request/response models
- **Handlers** in `handler/` — HTTP endpoint logic
- **Router** registration in `cmd/server/infra.go`

## ID Translation
Repositories expose both external UUID (`sub_id`) and internal `int32` ID:
```go
func (r *TransactionRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error)
func (r *TransactionRepository) GetBySubID(ctx context.Context, subID string) (*models.Transaction, error)
```

## Validation Pattern
Use `syncerr.Group` for parallel entity validation:
```go
grp := syncerr.Group{}
grp.Go(func() error {
    id, err := s.accountService.GetIDBySubID(ctx, accountID)
    if id == 0 {
        return fmt.Errorf("%w: %s", entities.ErrAccountNotFound, accountID)
    }
    ids.accountID = id
    return nil
})
if errs := grp.Wait(); len(errs) > 0 {
    return nil, errs[0]
}
```

## Error Wrapping
```go
fmt.Errorf("%w: %s", entities.ErrAccountNotFound, accountID)
fmt.Errorf("failed to create transaction from draft: %w", err)
```
