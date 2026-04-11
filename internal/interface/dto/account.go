package dto

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

func isValidAccountType(t string) bool {
	switch t {
	case "expense", "income", "transfer":
		return true
	default:
		return false
	}
}

func ParseAccountListParams(r *http.Request) *models.AccountSearchParams {
	_, span := observability.StartDTOSpan(r.Context(), "account", "parse_list_params")
	defer span.End()

	q := r.URL.Query()
	params := &models.AccountSearchParams{
		PageNumber: 1,
		PageSize:   10,
	}

	if v := q.Get("q"); v != "" {
		params.Q = &v
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

func ValidateCreateAccountRequest(req *models.CreateAccountRequest) error {
	_, span := observability.StartDTOSpan(context.Background(), "account", "validate_create")
	defer span.End()

	if req.Name == "" {
		return entities.ErrNameRequired
	}
	if req.Type == "" {
		return entities.ErrTypeRequired
	}
	if !isValidAccountType(string(req.Type)) {
		return entities.ErrInvalidType
	}
	if req.Balance < 0 {
		return entities.ErrNegativeBalance
	}
	return nil
}

func ValidateUpdateAccountRequest(req *models.UpdateAccountRequest) error {
	_, span := observability.StartDTOSpan(context.Background(), "account", "validate_update")
	defer span.End()

	if req.Type != nil && !isValidAccountType(string(*req.Type)) {
		return entities.ErrInvalidType
	}
	if req.Balance != nil && *req.Balance < 0 {
		return entities.ErrNegativeBalance
	}
	return nil
}
