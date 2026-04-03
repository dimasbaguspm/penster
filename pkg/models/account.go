package models

import (
	"time"

	"github.com/dimasbaguspm/penster/pkg/response"
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
	response.Response
	Data Account `json:"data"`
}

type AccountsResponse struct {
	response.PaginatedResponse
	Data []Account `json:"data"`
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
