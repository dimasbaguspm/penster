package dto

import "errors"

var (
	ErrNameRequired             = errors.New("name is required")
	ErrTypeRequired             = errors.New("type is required")
	ErrInvalidAccountType       = errors.New("invalid account type, must be one of: expense, income, transfer")
	ErrInvalidCategoryType      = errors.New("invalid category type, must be one of: expense, income, transfer")
	ErrAccountIDRequired        = errors.New("account_id is required")
	ErrCategoryIDRequired       = errors.New("category_id is required")
	ErrTransactionTypeRequired  = errors.New("transaction_type is required")
	ErrTitleRequired            = errors.New("title is required")
	ErrInvalidAmount            = errors.New("amount must be greater than 0")
	ErrCurrencyRequired         = errors.New("currency is required")
	ErrInvalidTransactionType   = errors.New("invalid transaction type, must be one of: expense, income, transfer")
)

func isValidAccountType(t string) bool {
	switch t {
	case "expense", "income", "transfer":
		return true
	default:
		return false
	}
}

func isValidCategoryType(t string) bool {
	switch t {
	case "expense", "income", "transfer":
		return true
	default:
		return false
	}
}

func isValidTransactionType(t string) bool {
	switch t {
	case "expense", "income", "transfer":
		return true
	default:
		return false
	}
}
