package models

import (
	"time"

	"github.com/dimasbaguspm/penster/pkg/response"
)

type DraftStatus string

const (
	DraftStatusPending   DraftStatus = "pending"
	DraftStatusConfirmed DraftStatus = "confirmed"
	DraftStatusRejected  DraftStatus = "rejected"
)

type DraftSource string

const (
	DraftSourceManual    DraftSource = "manual"
	DraftSourceIngestion DraftSource = "ingestion"
)

type Draft struct {
	ID                string     `json:"-"`
	SubID             string     `json:"id"`
	AccountID         string     `json:"account_id"`
	TransferAccountID string     `json:"transfer_account_id"`
	CategoryID        string     `json:"category_id"`
	TransactionType   string     `json:"transaction_type"`
	Title             string     `json:"title"`
	Amount            int64      `json:"amount"`
	Currency          string     `json:"currency"`
	CurrencyRate      float64    `json:"currency_rate"`
	TransactedAt      time.Time  `json:"transacted_at"`
	Notes             string     `json:"notes"`
	Source            string     `json:"source"`
	Status            string     `json:"status"`
	ConfirmedAt       *time.Time `json:"confirmed_at"`
	RejectedAt        *time.Time `json:"rejected_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type CreateDraftRequest struct {
	AccountID         string `json:"account_id" binding:"required"`
	TransferAccountID string `json:"transfer_account_id,omitempty"`
	CategoryID        string `json:"category_id" binding:"required"`
	TransactionType   string `json:"transaction_type" binding:"required,oneof=expense income transfer"`
	Title             string `json:"title" binding:"required"`
	Amount            int64  `json:"amount" binding:"required,gt=0"`
	Currency          string `json:"currency" binding:"required"`
	TransactedAt      string `json:"transacted_at" binding:"required"`
	Notes             string `json:"notes,omitempty"`
	Source            string `json:"source" binding:"required,oneof=manual ingestion"`
}

type UpdateDraftRequest struct {
	AccountID         *string `json:"account_id,omitempty"`
	TransferAccountID *string `json:"transfer_account_id,omitempty"`
	CategoryID        *string `json:"category_id,omitempty"`
	TransactionType   *string `json:"transaction_type,omitempty"`
	Title             *string `json:"title,omitempty"`
	Amount            *int64  `json:"amount,omitempty"`
	Currency          *string `json:"currency,omitempty"`
	TransactedAt      *string `json:"transacted_at,omitempty"`
	Notes             *string `json:"notes,omitempty"`
}

type DraftResponse struct {
	response.Response
	Data Draft `json:"data"`
}

type DraftsResponse struct {
	response.PaginatedResponse
	Data []Draft `json:"data"`
}

type DraftSearchParams struct {
	Source   *string
	Status   *string
	PageSize int
}

func (p *DraftSearchParams) Offset() int {
	if p.PageSize <= 0 {
		return 0
	}
	return p.PageSize
}
