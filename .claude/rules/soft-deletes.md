---
name: soft-deletes
description: Soft delete implementation pattern in Penster
type: rule
paths:
  - "internal/domain/entities/**/*.go"
  - "internal/domain/repository/**/*.go"
  - "internal/infrastructure/database/**/*.go"
  - "migrations/*.sql"
---

# Soft Deletes in Penster

## Affected Entities
- Account
- Category
- Transaction
- Draft

## Implementation
- `deleted_at` timestamp column (NULL = active, timestamp = deleted)
- Repository queries filter: `WHERE deleted_at IS NULL`
- No hard deletes — all major entities support soft delete

## Behavior
- Deleted entities are excluded from all queries
- FK relationships handle cascading soft deletes where needed
- Only `rejected` drafts can be deleted (enforced by DraftService)
