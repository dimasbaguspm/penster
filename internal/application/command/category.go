package command

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

// CategoryCommandInterface defines write operations for categories
type CategoryCommandInterface interface {
	Create(ctx context.Context, params query.CreateCategoryParams) (*models.Category, error)
	Update(ctx context.Context, id string, params query.UpdateCategoryParams) (*models.Category, error)
	Delete(ctx context.Context, id string) (*models.Category, error)
}

var _ CategoryCommandInterface = (*CategoryCommand)(nil)

// CategoryCommand implements CategoryCommandInterface using repository
type CategoryCommand struct {
	repo *repository.CategoryRepository
}

// NewCategoryCommand creates a new CategoryCommand
func NewCategoryCommand(repo *repository.CategoryRepository) *CategoryCommand {
	return &CategoryCommand{repo: repo}
}

func (c *CategoryCommand) Create(ctx context.Context, params query.CreateCategoryParams) (*models.Category, error) {
	ctx, span := observability.StartCommandSpan(ctx, "category", "create")
	defer span.End()
	return c.repo.Create(ctx, params)
}

func (c *CategoryCommand) Update(ctx context.Context, id string, params query.UpdateCategoryParams) (*models.Category, error) {
	ctx, span := observability.StartCommandSpan(ctx, "category", "update")
	defer span.End()
	return c.repo.UpdateBySubID(ctx, id, params)
}

func (c *CategoryCommand) Delete(ctx context.Context, id string) (*models.Category, error) {
	ctx, span := observability.StartCommandSpan(ctx, "category", "delete")
	defer span.End()
	return c.repo.DeleteBySubID(ctx, id)
}
