package entities

import "errors"

// Transaction-specific errors
var (
	ErrTransferAccountNotFound = errors.New("transfer account not found")
	ErrTransferToSameAccount   = errors.New("transfer_account_id cannot be the same as account_id")
	ErrTransactionTypeRequired = errors.New("transaction_type is required")
	ErrTitleRequired           = errors.New("title is required")
	ErrInvalidAmount           = errors.New("amount must be greater than 0")
	ErrCurrencyRequired        = errors.New("currency is required")
	ErrInvalidTransactionType  = errors.New("invalid transaction type, must be one of: expense, income, transfer")
)
