package entities

import "errors"

// Entity represents a domain entity with an ID
type Entity struct {
	ID string
}

// Generic validation errors
var (
	ErrNameRequired = errors.New("name is required")
	ErrTypeRequired = errors.New("type is required")
	ErrInvalidType  = errors.New("invalid type, must be one of: expense, income, transfer")
	ErrIDRequired   = errors.New("id is required")
	ErrInvalidID    = errors.New("id must be a valid UUID")
	ErrEmptyID      = errors.New("id cannot be empty")
)
