package models

import (
	"time"
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
	Data Category `json:"data"`
}

type CategoryPagedResponse struct {
	Items      []Category `json:"items"`
	PageSize   int        `json:"page_size"`
	PageNumber int        `json:"page_number"`
	TotalPages int        `json:"total_pages"`
	TotalItems int64      `json:"total_items"`
}

func NewCategoryPagedResponse(items []Category, pageSize, pageNumber int, totalItems int64) CategoryPagedResponse {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}
	return CategoryPagedResponse{
		Items:      items,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
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
