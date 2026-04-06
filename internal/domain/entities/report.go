package entities

import "errors"

// Report-specific errors
var (
	ErrInvalidDateRange = errors.New("invalid date range: start_date must be before end_date")
)
