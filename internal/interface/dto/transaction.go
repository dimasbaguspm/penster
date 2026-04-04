package dto

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/pkg/models"
)

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func isValidTransactionType(t string) bool {
	switch t {
	case "expense", "income", "transfer":
		return true
	default:
		return false
	}
}

func ParseTransactionListParams(r *http.Request) *models.TransactionSearchParams {
	q := r.URL.Query()
	params := &models.TransactionSearchParams{
		PageNumber: 1,
		PageSize:   10,
	}

	if v := q.Get("q"); v != "" {
		params.Q = &v
	}
	if v := q.Get("account_id"); v != "" {
		params.AccountIDs = []string{v}
	}
	if v := q.Get("category_id"); v != "" {
		params.CategoryIDs = []string{v}
	}
	if v := q.Get("transaction_type"); v != "" {
		params.TransactionTypes = []string{v}
	}
	if v := q.Get("sort_by"); v != "" {
		params.SortBy = v
	}
	if v := q.Get("sort_order"); v != "" {
		params.SortOrder = v
	}
	if v := q.Get("page"); v != "" {
		if page, err := strconv.Atoi(v); err == nil {
			params.PageNumber = page
		}
	}
	if v := q.Get("page_size"); v != "" {
		if pageSize, err := strconv.Atoi(v); err == nil {
			params.PageSize = pageSize
		}
	}

	return params
}

func ValidateCreateTransactionRequest(req *models.CreateTransactionRequest) error {
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
	if req.Title == "" {
		return entities.ErrTitleRequired
	}
	if req.Amount <= 0 {
		return entities.ErrInvalidAmount
	}
	if req.Currency == "" {
		return entities.ErrCurrencyRequired
	}
	if !isValidTransactionType(string(req.TransactionType)) {
		return entities.ErrInvalidTransactionType
	}
	if req.TransferAccountID != "" && !isValidUUID(req.TransferAccountID) {
		return entities.ErrInvalidID
	}
	if req.TransactionType == models.TransactionTypeTransfer && req.TransferAccountID == req.AccountID {
		return entities.ErrTransferToSameAccount
	}
	return nil
}

func ValidateUpdateTransactionRequest(req *models.UpdateTransactionRequest) error {
	if req.TransactionType != nil && !isValidTransactionType(string(*req.TransactionType)) {
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
