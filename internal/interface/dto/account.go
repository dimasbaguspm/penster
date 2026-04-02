package dto

import (
	"net/http"
	"strconv"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func ParseAccountListParams(r *http.Request) *models.AccountSearchParams {
	q := r.URL.Query()
	params := &models.AccountSearchParams{}

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
	if req.Name == "" {
		return ErrNameRequired
	}
	if req.Type == "" {
		return ErrTypeRequired
	}
	if !isValidAccountType(string(req.Type)) {
		return ErrInvalidAccountType
	}
	return nil
}

func ValidateUpdateAccountRequest(req *models.UpdateAccountRequest) error {
	if req.Type != nil && !isValidAccountType(string(*req.Type)) {
		return ErrInvalidAccountType
	}
	return nil
}
