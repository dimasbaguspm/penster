package models

import (
	"time"

	"github.com/dimasbaguspm/penster/pkg/response"
)

type CategoryType string

const (
	CategoryTypeExpense  CategoryType = "expense"
	CategoryTypeIncome   CategoryType = "income"
	CategoryTypeTransfer CategoryType = "transfer"
)

type Category struct {
	ID        string       `json:"-"`
	SubID     string       `json:"id"`
	Name      string       `json:"name"`
	Type      CategoryType `json:"type"`
	DeletedAt *time.Time   `json:"deleted_at,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name string       `json:"name" binding:"required"`
	Type CategoryType `json:"type" binding:"required,oneof=expense income transfer"`
}

type UpdateCategoryRequest struct {
	Name *string       `json:"name,omitempty"`
	Type *CategoryType `json:"type,omitempty"`
}

type CategoryResponse struct {
	response.Response
	Data Category `json:"data"`
}

type CategoriesResponse struct {
	response.PaginatedResponse
	Data []Category `json:"data"`
}

type CategorySearchParams struct {
	SubID      *string
	Q          *string
	SortBy     string
	SortOrder  string
	PageSize   int
	PageNumber int
}

func (p *CategorySearchParams) Offset() int {
	if p.PageSize <= 0 {
		return 0
	}
	return (p.PageNumber - 1) * p.PageSize
}
