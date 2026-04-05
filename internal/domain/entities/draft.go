package entities

import "errors"

var (
	ErrDraftNotFound         = errors.New("draft not found")
	ErrDraftNotPending       = errors.New("draft is not in pending status")
	ErrDraftAlreadyConfirmed = errors.New("draft is already confirmed")
	ErrDraftAlreadyRejected  = errors.New("draft is already rejected")
	ErrDraftNotRejected      = errors.New("draft must be rejected before deletion")
	ErrSourceRequired        = errors.New("source is required")
	ErrInvalidDraftSource    = errors.New("invalid source, must be one of: manual, ingestion")
)
