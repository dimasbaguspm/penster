package query

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

// CategoryQueryInterface defines read operations for categories
type CategoryQueryInterface interface {
	GetByID(ctx context.Context, id string) (*models.Category, error)
	GetIDBySubID(ctx context.Context, subID string) (int32, error)
	List(ctx context.Context, params *models.CategorySearchParams) ([]*models.Category, int64, error)
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
	return q.repo.GetBySubID(ctx, id)
}

func (q *CategoryQuery) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	return q.repo.GetIDBySubID(ctx, subID)
}

func (q *CategoryQuery) List(ctx context.Context, params *models.CategorySearchParams) ([]*models.Category, int64, error) {
	return q.repo.List(ctx, params)
}
