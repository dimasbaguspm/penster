package command

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

// CategoryCommandInterface defines write operations for categories
type CategoryCommandInterface interface {
	Create(ctx context.Context, req *models.CreateCategoryRequest) (*models.Category, error)
	Update(ctx context.Context, id string, req *models.UpdateCategoryRequest) (*models.Category, error)
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

func (c *CategoryCommand) Create(ctx context.Context, req *models.CreateCategoryRequest) (*models.Category, error) {
	return c.repo.Create(ctx, req)
}

func (c *CategoryCommand) Update(ctx context.Context, id string, req *models.UpdateCategoryRequest) (*models.Category, error) {
	return c.repo.UpdateBySubID(ctx, id, req)
}

func (c *CategoryCommand) Delete(ctx context.Context, id string) (*models.Category, error) {
	return c.repo.DeleteBySubID(ctx, id)
}
