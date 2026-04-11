package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type CategoryRepository struct {
	db *query.Queries
}

func NewCategoryRepository(db *query.Queries) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, params query.CreateCategoryParams) (*models.Category, error) {
	ctx, span := observability.StartRepoSpan(ctx, "categories", "Create")
	defer span.End()

	result, err := r.db.CreateCategory(ctx, params)
	if err != nil {
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toCategoryModel(ctx, result), nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int32) (*models.Category, error) {
	ctx, span := observability.StartRepoSpan(ctx, "categories", "GetByID")
	defer span.End()

	result, err := r.db.GetCategoryByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toCategoryModel(ctx, result), nil
}

func (r *CategoryRepository) GetBySubID(ctx context.Context, subID string) (*models.Category, error) {
	ctx, span := observability.StartRepoSpan(ctx, "categories", "GetBySubID")
	defer span.End()

	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetCategoryBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toCategoryModel(ctx, result), nil
}

func (r *CategoryRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	ctx, span := observability.StartRepoSpan(ctx, "categories", "GetIDBySubID")
	defer span.End()

	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetCategoryBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		observability.RecordError(ctx, err)
		return 0, err
	}
	return result.ID, nil
}

func (r *CategoryRepository) UpdateBySubID(ctx context.Context, subID string, params query.UpdateCategoryParams) (*models.Category, error) {
	ctx, span := observability.StartRepoSpan(ctx, "categories", "UpdateBySubID")
	defer span.End()

	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Update(ctx, id, params)
}

func (r *CategoryRepository) DeleteBySubID(ctx context.Context, subID string) (*models.Category, error) {
	ctx, span := observability.StartRepoSpan(ctx, "categories", "DeleteBySubID")
	defer span.End()

	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Delete(ctx, id)
}

func (r *CategoryRepository) List(ctx context.Context, params query.ListCategoriesParams) ([]*models.Category, int64, error) {
	ctx, span := observability.StartRepoSpan(ctx, "categories", "List")
	defer span.End()

	rows, err := r.db.ListCategories(ctx, params)
	if err != nil {
		observability.RecordError(ctx, err)
		return nil, 0, err
	}

	categories := make([]*models.Category, 0, len(rows))
	var total int64
	for _, row := range rows {
		categories = append(categories, toCategoryModel(ctx, query.Category{
			ID:        row.ID,
			SubID:     row.SubID,
			Name:      row.Name,
			Type:      row.Type,
			DeletedAt: row.DeletedAt,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}))
		total = row.Total
	}

	return categories, total, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id int32, params query.UpdateCategoryParams) (*models.Category, error) {
	ctx, span := observability.StartRepoSpan(ctx, "categories", "Update")
	defer span.End()

	result, err := r.db.UpdateCategory(ctx, query.UpdateCategoryParams{
		Name: params.Name,
		Type: params.Type,
		ID:   id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toCategoryModel(ctx, result), nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int32) (*models.Category, error) {
	ctx, span := observability.StartRepoSpan(ctx, "categories", "Delete")
	defer span.End()

	result, err := r.db.DeleteCategory(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toCategoryModel(ctx, result), nil
}

func toCategoryModel(ctx context.Context, q query.Category) *models.Category {
	_, span := observability.StartRepoSpan(ctx, "categories", "to_model")
	defer span.End()

	m := &models.Category{
		SubID:     uuid.UUID(q.SubID.Bytes).String(),
		Name:      q.Name,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	m.Type = models.CategoryType(q.Type)

	if q.DeletedAt.Valid {
		m.DeletedAt = &q.DeletedAt.Time
	}
	if q.CreatedAt.Valid {
		m.CreatedAt = q.CreatedAt.Time
	}
	if q.UpdatedAt.Valid {
		m.UpdatedAt = q.UpdatedAt.Time
	}

	return m
}
