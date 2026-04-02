package service

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/application/command"
	"github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type CategoryService struct {
	query    query.CategoryQueryInterface
	commands command.CategoryCommandInterface
}

func NewCategoryService(query query.CategoryQueryInterface, commands command.CategoryCommandInterface) *CategoryService {
	return &CategoryService{
		query:    query,
		commands: commands,
	}
}

func (s *CategoryService) Create(ctx context.Context, req *models.CreateCategoryRequest) (*models.Category, error) {
	return s.commands.Create(ctx, req)
}

func (s *CategoryService) GetByID(ctx context.Context, id string) (*models.Category, error) {
	return s.query.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context, params *models.CategorySearchParams) ([]*models.Category, int64, error) {
	return s.query.List(ctx, params)
}

func (s *CategoryService) Update(ctx context.Context, id string, req *models.UpdateCategoryRequest) (*models.Category, error) {
	return s.commands.Update(ctx, id, req)
}

func (s *CategoryService) Delete(ctx context.Context, id string) (*models.Category, error) {
	return s.commands.Delete(ctx, id)
}
