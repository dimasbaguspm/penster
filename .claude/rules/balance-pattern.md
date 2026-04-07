---
name: balance-pattern
description: Transaction balance update/reverse conventions in Penster
type: rule
paths:
  - "internal/application/service/**/*.go"
  - "internal/application/command/**/*.go"
  - "internal/domain/entities/**/*.go"
---

# Balance Update Pattern in Penster

## Core Rule
**Always reverse first, then apply new values.**

## Transaction Types and Balance Effects

| Type | Effect |
|------|--------|
| `expense` | `account.balance -= amount` |
| `income` | `account.balance += amount` |
| `transfer` | `source.balance -= amount`, `dest.balance += amount` |

## Update Flow
1. **Reverse** existing balance change (undo old values)
2. **Apply** new balance change (apply new values)

## Delete Flow
1. **Reverse** the balance change (undo)

## Create Flow
1. **Apply** balance change directly (no reversal needed)

## Implementation
See `AccountService.UpdateAccountBalances()` and `AccountService.ReverseAccountBalances()`.

## Same-Account Transfer Prevention
Transfers where `account_id == transfer_account_id` must be rejected with `ErrTransferToSameAccount`.

## Enhanced Amount
`enhanced_amount = base_amount * currency_rate`
This stores the amount in base currency (e.g., EUR) for cross-currency comparison.
