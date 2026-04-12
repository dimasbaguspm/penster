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
	log := observability.NewLogger(ctx, "repository", "category")
	ctx, span := observability.StartRepoSpan(log.Context(), "category", "Create")
	defer span.End()

	log.Info("category.create started", "name", params.Name)
	result, err := r.db.CreateCategory(ctx, params)
	if err != nil {
		log.Error("category.create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("category.created", "id", result.ID)
	return toCategoryModel(ctx, result), nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int32) (*models.Category, error) {
	log := observability.NewLogger(ctx, "repository", "category")
	ctx, span := observability.StartRepoSpan(log.Context(), "category", "GetByID")
	defer span.End()

	log.Info("category.get_by_id started", "id", id)
	result, err := r.db.GetCategoryByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("category.get_by_id not found", "id", id)
			return nil, nil
		}
		log.Error("category.get_by_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("category.get_by_id succeeded", "id", id)
	return toCategoryModel(ctx, result), nil
}

func (r *CategoryRepository) GetBySubID(ctx context.Context, subID string) (*models.Category, error) {
	log := observability.NewLogger(ctx, "repository", "category")
	ctx, span := observability.StartRepoSpan(log.Context(), "category", "GetBySubID")
	defer span.End()

	log.Info("category.get_by_sub_id started", "sub_id", subID)
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetCategoryBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("category.get_by_sub_id not found", "sub_id", subID)
			return nil, nil
		}
		log.Error("category.get_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("category.get_by_sub_id succeeded", "sub_id", subID)
	return toCategoryModel(ctx, result), nil
}

func (r *CategoryRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	log := observability.NewLogger(ctx, "repository", "category")
	ctx, span := observability.StartRepoSpan(log.Context(), "category", "GetIDBySubID")
	defer span.End()

	log.Info("category.get_id_by_sub_id started", "sub_id", subID)
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetCategoryBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("category.get_id_by_sub_id not found", "sub_id", subID)
			return 0, nil
		}
		log.Error("category.get_id_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return 0, err
	}
	log.Info("category.get_id_by_sub_id succeeded", "sub_id", subID, "id", result.ID)
	return result.ID, nil
}

func (r *CategoryRepository) UpdateBySubID(ctx context.Context, subID string, params query.UpdateCategoryParams) (*models.Category, error) {
	log := observability.NewLogger(ctx, "repository", "category")
	ctx, span := observability.StartRepoSpan(log.Context(), "category", "UpdateBySubID")
	defer span.End()

	log.Info("category.update_by_sub_id started", "sub_id", subID)
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		log.Debug("category.update_by_sub_id not found", "sub_id", subID)
		return nil, nil
	}
	return r.Update(ctx, id, params)
}

func (r *CategoryRepository) DeleteBySubID(ctx context.Context, subID string) (*models.Category, error) {
	log := observability.NewLogger(ctx, "repository", "category")
	ctx, span := observability.StartRepoSpan(log.Context(), "category", "DeleteBySubID")
	defer span.End()

	log.Info("category.delete_by_sub_id started", "sub_id", subID)
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		log.Debug("category.delete_by_sub_id not found", "sub_id", subID)
		return nil, nil
	}
	return r.Delete(ctx, id)
}

func (r *CategoryRepository) List(ctx context.Context, params query.ListCategoriesParams) ([]*models.Category, int64, error) {
	log := observability.NewLogger(ctx, "repository", "category")
	ctx, span := observability.StartRepoSpan(log.Context(), "category", "List")
	defer span.End()

	log.Info("category.list started")
	rows, err := r.db.ListCategories(ctx, params)
	if err != nil {
		log.Error("category.list failed", "error", err)
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

	log.Info("category.list succeeded", "count", len(categories), "total", total)
	return categories, total, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id int32, params query.UpdateCategoryParams) (*models.Category, error) {
	log := observability.NewLogger(ctx, "repository", "category")
	ctx, span := observability.StartRepoSpan(log.Context(), "category", "Update")
	defer span.End()

	log.Info("category.update started", "id", id)
	result, err := r.db.UpdateCategory(ctx, query.UpdateCategoryParams{
		Name: params.Name,
		Type: params.Type,
		ID:   id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("category.update not found", "id", id)
			return nil, nil
		}
		log.Error("category.update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("category.updated", "id", id)
	return toCategoryModel(ctx, result), nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int32) (*models.Category, error) {
	log := observability.NewLogger(ctx, "repository", "category")
	ctx, span := observability.StartRepoSpan(log.Context(), "category", "Delete")
	defer span.End()

	log.Info("category.delete started", "id", id)
	result, err := r.db.DeleteCategory(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("category.delete not found", "id", id)
			return nil, nil
		}
		log.Error("category.delete failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("category.deleted", "id", id)
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
