package dto

import "errors"

var (
	ErrNameRequired       = errors.New("name is required")
	ErrTypeRequired      = errors.New("type is required")
	ErrInvalidAccountType  = errors.New("invalid account type, must be one of: expense, income, transfer")
	ErrInvalidCategoryType = errors.New("invalid category type, must be one of: expense, income, transfer")
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
