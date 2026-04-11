package valueobjects

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
	"github.com/jackc/pgx/v5/pgtype"
)

// ToCreateCategoryParams converts CreateCategoryRequest to query params
func ToCreateCategoryParams(ctx context.Context, req *models.CreateCategoryRequest) query.CreateCategoryParams {
	_, span := observability.StartValueObjectSpan(ctx, "category", "to_create_params")
	defer span.End()

	return query.CreateCategoryParams{
		Name: req.Name,
		Type: string(req.Type),
	}
}

// ToUpdateCategoryParams converts UpdateCategoryRequest to query params
func ToUpdateCategoryParams(ctx context.Context, req *models.UpdateCategoryRequest) query.UpdateCategoryParams {
	_, span := observability.StartValueObjectSpan(ctx, "category", "to_update_params")
	defer span.End()

	name := ""
	if req.Name != nil {
		name = *req.Name
	}
	catType := ""
	if req.Type != nil {
		catType = string(*req.Type)
	}

	return query.UpdateCategoryParams{
		Name: name,
		Type: catType,
	}
}

// ToListCategoriesParams converts CategorySearchParams to query params
func ToListCategoriesParams(ctx context.Context, params *models.CategorySearchParams) query.ListCategoriesParams {
	_, span := observability.StartValueObjectSpan(ctx, "category", "to_list_params")
	defer span.End()

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
