package entities

import (
	"errors"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/jackc/pgx/v5/pgtype"
)

// Category-specific errors
var (
	ErrCategoryNotFound = errors.New("category not found")
)

// ToListCategoriesParams converts CategorySearchParams to query params
func ToListCategoriesParams(params *models.CategorySearchParams) query.ListCategoriesParams {
	var subID pgtype.UUID
	if params.SubID != nil {
		subID = pgtype.UUID{Bytes: conv.ParseUUID(*params.SubID), Valid: true}
	}

	return query.ListCategoriesParams{
		Column1: subID,
		Column2: conv.StringPtrToEmpty(params.Q),
		Column3: params.SortBy,
		Column4: params.SortOrder,
		Column5: "",
		Column6: pgtype.Timestamptz{},
		Column7: pgtype.Timestamptz{},
		Column8: int64(params.PageSize),
	}
}
