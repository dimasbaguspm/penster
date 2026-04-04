package repository

import "errors"

var (
	ErrAccountNotFound         = errors.New("account not found")
	ErrCategoryNotFound        = errors.New("category not found")
	ErrTransferAccountNotFound = errors.New("transfer account not found")
	ErrTransferToSameAccount   = errors.New("transfer to same account is not allowed")
)
