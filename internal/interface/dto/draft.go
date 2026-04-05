package dto

import (
	"net/http"
	"strconv"

	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/pkg/models"
)

func isValidDraftSource(s string) bool {
	switch s {
	case "manual", "ingestion":
		return true
	default:
		return false
	}
}

func ParseDraftListParams(r *http.Request) *models.DraftSearchParams {
	q := r.URL.Query()
	params := &models.DraftSearchParams{
		PageSize: 10,
	}

	if v := q.Get("source"); v != "" {
		params.Source = &v
	}
	if v := q.Get("status"); v != "" {
		params.Status = &v
	}
	if v := q.Get("page_size"); v != "" {
		if pageSize, err := strconv.Atoi(v); err == nil {
			params.PageSize = pageSize
		}
	}

	return params
}

func ValidateCreateDraftRequest(req *models.CreateDraftRequest) error {
	if req.AccountID == "" {
		return entities.ErrIDRequired
	}
	if !isValidUUID(req.AccountID) {
		return entities.ErrInvalidID
	}
	if req.CategoryID == "" {
		return entities.ErrIDRequired
	}
	if !isValidUUID(req.CategoryID) {
		return entities.ErrInvalidID
	}
	if req.TransactionType == "" {
		return entities.ErrTransactionTypeRequired
	}
	if !isValidTransactionType(req.TransactionType) {
		return entities.ErrInvalidTransactionType
	}
	if req.Title == "" {
		return entities.ErrTitleRequired
	}
	if req.Amount <= 0 {
		return entities.ErrInvalidAmount
	}
	if req.Currency == "" {
		return entities.ErrCurrencyRequired
	}
	if req.TransactedAt == "" {
		return entities.ErrTransactedAtRequired
	}
	if req.Source == "" {
		return entities.ErrSourceRequired
	}
	if !isValidDraftSource(req.Source) {
		return entities.ErrInvalidDraftSource
	}
	if req.TransactionType == string(models.TransactionTypeTransfer) && req.TransferAccountID != "" && req.AccountID == req.TransferAccountID {
		return entities.ErrTransferToSameAccount
	}
	if req.TransferAccountID != "" && !isValidUUID(req.TransferAccountID) {
		return entities.ErrInvalidID
	}
	return nil
}

func ValidateUpdateDraftRequest(req *models.UpdateDraftRequest) error {
	if req.TransactionType != nil && !isValidTransactionType(*req.TransactionType) {
		return entities.ErrInvalidTransactionType
	}
	if req.Amount != nil && *req.Amount <= 0 {
		return entities.ErrInvalidAmount
	}
	if req.AccountID != nil && *req.AccountID == "" {
		return entities.ErrEmptyID
	}
	if req.AccountID != nil && !isValidUUID(*req.AccountID) {
		return entities.ErrInvalidID
	}
	if req.CategoryID != nil && *req.CategoryID == "" {
		return entities.ErrEmptyID
	}
	if req.CategoryID != nil && !isValidUUID(*req.CategoryID) {
		return entities.ErrInvalidID
	}
	if req.TransferAccountID != nil && *req.TransferAccountID == "" {
		return entities.ErrEmptyID
	}
	if req.TransferAccountID != nil && !isValidUUID(*req.TransferAccountID) {
		return entities.ErrInvalidID
	}
	return nil
}
