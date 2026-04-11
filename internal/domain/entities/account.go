package entities

import "errors"

// Account-specific errors
var (
	ErrAccountNotFound = errors.New("account not found")
)
