package dto

import (
	"net/http"
	"strconv"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func ParseCategoryListParams(r *http.Request) *models.CategorySearchParams {
	q := r.URL.Query()
	params := &models.CategorySearchParams{
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

func ValidateCreateCategoryRequest(req *models.CreateCategoryRequest) error {
	if req.Name == "" {
		return ErrNameRequired
	}
	if req.Type == "" {
		return ErrTypeRequired
	}
	if !isValidCategoryType(string(req.Type)) {
		return ErrInvalidCategoryType
	}
	return nil
}

func ValidateUpdateCategoryRequest(req *models.UpdateCategoryRequest) error {
	if req.Type != nil && !isValidCategoryType(string(*req.Type)) {
		return ErrInvalidCategoryType
	}
	return nil
}
