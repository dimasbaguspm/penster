package dto

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
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
	_, span := observability.StartDTOSpan(r.Context(), "category", "parse_list_params")
	defer span.End()

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

func ValidateCreateCategoryRequest(ctx context.Context, req *models.CreateCategoryRequest) error {
	_, span := observability.StartDTOSpan(ctx, "category", "validate_create")
	defer span.End()

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

func ValidateUpdateCategoryRequest(ctx context.Context, req *models.UpdateCategoryRequest) error {
	_, span := observability.StartDTOSpan(ctx, "category", "validate_update")
	defer span.End()

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return entities.ErrNameRequired
	}
	if req.Type != nil && !isValidCategoryType(string(*req.Type)) {
		return entities.ErrInvalidType
	}
	return nil
}
