---
name: draft-lifecycle
description: Draft lifecycle state machine rules in Penster
type: rule
paths:
  - "internal/application/service/draft.go"
  - "internal/application/command/draft.go"
  - "internal/domain/entities/draft.go"
  - "internal/interface/handler/draft.go"
---

# Draft Lifecycle State Machine

## States
- `pending` — awaiting confirmation or rejection
- `confirmed` — becomes a Transaction (irreversible)
- `rejected` — discarded (can be deleted)

## Transition Rules

| Action | Allowed From States |
|--------|-------------------|
| **Confirm** | `pending` only |
| **Reject** | `pending` only |
| **Delete** | `rejected` only |

## Confirm Side Effects
When a draft is confirmed:
1. Creates a `Transaction` record
2. Updates account balances (via `AccountService.UpdateAccountBalances()`)

## Sources
- `manual` — user-created via API
- `ingestion` — imported data

## Enforcement
Rules are enforced in `DraftService` methods:
- `ConfirmDraft`: only `pending` drafts
- `RejectDraft`: only `pending` drafts
- `DeleteDraft`: only `rejected` drafts
