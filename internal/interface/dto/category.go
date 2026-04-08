package dto

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/pkg/models"
)

func isValidCategoryType(t string) bool {
	switch t {
	case "expense", "income", "transfer":
		return true
	default:
		return false
	}
}

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
	if strings.TrimSpace(req.Name) == "" {
		return entities.ErrNameRequired
	}
	if req.Type == "" {
		return entities.ErrTypeRequired
	}
	if !isValidCategoryType(string(req.Type)) {
		return entities.ErrInvalidType
	}
	return nil
}

func ValidateUpdateCategoryRequest(req *models.UpdateCategoryRequest) error {
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return entities.ErrNameRequired
	}
	if req.Type != nil && !isValidCategoryType(string(*req.Type)) {
		return entities.ErrInvalidType
	}
	return nil
}
