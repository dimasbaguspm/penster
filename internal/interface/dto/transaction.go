package dto

import (
	"net/http"
	"strconv"

	"github.com/dimasbaguspm/penster/pkg/models"
)

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
		return ErrAccountIDRequired
	}
	if req.CategoryID == "" {
		return ErrCategoryIDRequired
	}
	if req.TransactionType == "" {
		return ErrTransactionTypeRequired
	}
	if req.Title == "" {
		return ErrTitleRequired
	}
	if req.Amount <= 0 {
		return ErrInvalidAmount
	}
	if req.Currency == "" {
		return ErrCurrencyRequired
	}
	if !isValidTransactionType(string(req.TransactionType)) {
		return ErrInvalidTransactionType
	}
	return nil
}

func ValidateUpdateTransactionRequest(req *models.UpdateTransactionRequest) error {
	if req.TransactionType != nil && !isValidTransactionType(string(*req.TransactionType)) {
		return ErrInvalidTransactionType
	}
	return nil
}
