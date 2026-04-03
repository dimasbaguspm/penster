package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type CategoryRepository struct {
	db *query.Queries
}

func NewCategoryRepository(db *query.Queries) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, req *models.CreateCategoryRequest) (*models.Category, error) {
	result, err := r.db.CreateCategory(ctx, query.CreateCategoryParams{
		Name: req.Name,
		Type: string(req.Type),
	})
	if err != nil {
		return nil, err
	}
	return toCategoryModel(result), nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int32) (*models.Category, error) {
	result, err := r.db.GetCategoryByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toCategoryModel(result), nil
}

func (r *CategoryRepository) GetBySubID(ctx context.Context, subID string) (*models.Category, error) {
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetCategoryBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toCategoryModel(result), nil
}

func (r *CategoryRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetCategoryBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return result.ID, nil
}

func (r *CategoryRepository) UpdateBySubID(ctx context.Context, subID string, req *models.UpdateCategoryRequest) (*models.Category, error) {
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Update(ctx, id, req)
}

func (r *CategoryRepository) DeleteBySubID(ctx context.Context, subID string) (*models.Category, error) {
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Delete(ctx, id)
}

func (r *CategoryRepository) List(ctx context.Context, params *models.CategorySearchParams) ([]*models.Category, int64, error) {
	queryParams := entities.ToListCategoriesParams(params)
	rows, err := r.db.ListCategories(ctx, queryParams)
	if err != nil {
		return nil, 0, err
	}

	categories := make([]*models.Category, 0, len(rows))
	var total int64
	for _, row := range rows {
		categories = append(categories, toCategoryModel(query.Category{
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

func (r *CategoryRepository) Update(ctx context.Context, id int32, req *models.UpdateCategoryRequest) (*models.Category, error) {
	name := ""
	if req.Name != nil {
		name = *req.Name
	}
	catType := ""
	if req.Type != nil {
		catType = string(*req.Type)
	}

	result, err := r.db.UpdateCategory(ctx, query.UpdateCategoryParams{
		Name: name,
		Type: catType,
		ID:   id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toCategoryModel(result), nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int32) (*models.Category, error) {
	result, err := r.db.DeleteCategory(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toCategoryModel(result), nil
}

func toCategoryModel(q query.Category) *models.Category {
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
