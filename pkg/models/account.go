package models

import (
	"time"
)

type AccountType string

const (
	AccountTypeExpense  AccountType = "expense"
	AccountTypeIncome   AccountType = "income"
	AccountTypeTransfer AccountType = "transfer"
)

type Account struct {
	ID        string      `json:"-"`
	SubID     string      `json:"id"`
	Name      string      `json:"name"`
	Type      AccountType `json:"type"`
	Balance   int64       `json:"balance"`
	DeletedAt *time.Time  `json:"deleted_at,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CreateAccountRequest struct {
	Name    string      `json:"name" binding:"required"`
	Type    AccountType `json:"type" binding:"required,oneof=expense income transfer"`
	Balance int64       `json:"balance"`
}

type UpdateAccountRequest struct {
	Name    *string      `json:"name,omitempty"`
	Type    *AccountType `json:"type,omitempty"`
	Balance *int64       `json:"balance,omitempty"`
}

type AccountResponse struct {
	Data Account `json:"data"`
}

type AccountPagedResponse struct {
	Items      []Account `json:"items"`
	PageSize   int       `json:"page_size"`
	PageNumber int       `json:"page_number"`
	TotalPages int       `json:"total_pages"`
	TotalItems int64     `json:"total_items"`
}

func NewAccountPagedResponse(items []Account, pageSize, pageNumber int, totalItems int64) AccountPagedResponse {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}
	return AccountPagedResponse{
		Items:      items,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
}

type AccountSearchParams struct {
	SubID      *string
	Q          *string
	SortBy     string
	SortOrder  string
	PageSize   int
	PageNumber int
}

func (p *AccountSearchParams) Offset() int {
	if p.PageSize <= 0 {
		return 0
	}
	return (p.PageNumber - 1) * p.PageSize
}
