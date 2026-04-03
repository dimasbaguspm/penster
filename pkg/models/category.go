package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
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

func (p *CategorySearchParams) ToQueryParams() query.ListCategoriesParams {
	var subID pgtype.UUID
	if p.SubID != nil {
		subID = pgtype.UUID{Bytes: conv.ParseUUID(*p.SubID), Valid: true}
	}

	return query.ListCategoriesParams{
		Column1: subID,
		Column2: conv.StringPtrToEmpty(p.Q),
		Column3: p.SortBy,
		Column4: p.SortOrder,
		Column5: "",
		Column6: pgtype.Timestamptz{},
		Column7: pgtype.Timestamptz{},
		Column8: int64(p.PageSize),
	}
}
