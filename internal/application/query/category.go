package query

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

// CategoryQueryInterface defines read operations for categories
type CategoryQueryInterface interface {
	GetByID(ctx context.Context, id string) (*models.Category, error)
	GetIDBySubID(ctx context.Context, subID string) (int32, error)
	List(ctx context.Context, params query.ListCategoriesParams) ([]*models.Category, int64, error)
}

var _ CategoryQueryInterface = (*CategoryQuery)(nil)

// CategoryQuery implements CategoryQueryInterface using repository
type CategoryQuery struct {
	repo *repository.CategoryRepository
}

// NewCategoryQuery creates a new CategoryQuery
func NewCategoryQuery(repo *repository.CategoryRepository) *CategoryQuery {
	return &CategoryQuery{repo: repo}
}

func (q *CategoryQuery) GetByID(ctx context.Context, id string) (*models.Category, error) {
	ctx, span := observability.StartQuerySpan(ctx, "category", "get_by_id")
	defer span.End()
	return q.repo.GetBySubID(ctx, id)
}

func (q *CategoryQuery) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	ctx, span := observability.StartQuerySpan(ctx, "category", "get_id_by_sub_id")
	defer span.End()
	return q.repo.GetIDBySubID(ctx, subID)
}

func (q *CategoryQuery) List(ctx context.Context, params query.ListCategoriesParams) ([]*models.Category, int64, error) {
	ctx, span := observability.StartQuerySpan(ctx, "category", "list")
	defer span.End()
	return q.repo.List(ctx, params)
}
