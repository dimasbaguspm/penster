package models

import (
	"time"

	"github.com/dimasbaguspm/penster/pkg/response"
)

type TransactionType string

const (
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Transaction struct {
	ID                string          `json:"-"`
	SubID             string          `json:"id"`
	AccountID         string          `json:"account_id"`
	TransferAccountID *string         `json:"transfer_account_id,omitempty"`
	CategoryID        *string         `json:"category_id,omitempty"`
	TransactionType   TransactionType `json:"transaction_type"`
	Title             string          `json:"title"`
	Amount            int64           `json:"amount"`
	Currency          string          `json:"currency"`
	Notes             *string         `json:"notes,omitempty"`
	DeletedAt         *time.Time      `json:"deleted_at,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type CreateTransactionRequest struct {
	AccountID         string          `json:"account_id" binding:"required"`
	TransferAccountID *string         `json:"transfer_account_id,omitempty"`
	CategoryID        string          `json:"category_id" binding:"required"`
	TransactionType   TransactionType `json:"transaction_type" binding:"required,oneof=expense income transfer"`
	Title             string          `json:"title" binding:"required"`
	Amount            int64           `json:"amount" binding:"required"`
	Currency          string          `json:"currency" binding:"required"`
	Notes             *string         `json:"notes,omitempty"`
}

type UpdateTransactionRequest struct {
	AccountID         *string          `json:"account_id,omitempty"`
	TransferAccountID *string          `json:"transfer_account_id,omitempty"`
	CategoryID        *string          `json:"category_id,omitempty"`
	TransactionType   *TransactionType `json:"transaction_type,omitempty"`
	Title             *string          `json:"title,omitempty"`
	Amount            *int64           `json:"amount,omitempty"`
	Currency          *string          `json:"currency,omitempty"`
	Notes             *string          `json:"notes,omitempty"`
}

type TransactionResponse struct {
	response.Response
	Data Transaction `json:"data"`
}

type TransactionsResponse struct {
	response.PaginatedResponse
	Data []Transaction `json:"data"`
}

type TransactionSearchParams struct {
	SubIDs           []string
	AccountIDs       []string
	CategoryIDs      []string
	TransactionTypes []string
	TransactedAt     *time.Time
	Q                *string
	SortBy           string
	SortOrder        string
	PageSize         int
	PageNumber       int
}

func (p *TransactionSearchParams) Offset() int {
	if p.PageSize <= 0 {
		return 0
	}
	return (p.PageNumber - 1) * p.PageSize
}
