package entities

import "errors"

// Category-specific errors
var (
	ErrCategoryNotFound = errors.New("category not found")
)
