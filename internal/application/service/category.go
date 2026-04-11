package service

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/application/command"
	"github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/domain/valueobjects"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
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
	ctx, span := observability.StartServiceSpan(ctx, "CategoryService", "Create")
	defer span.End()

	params := valueobjects.ToCreateCategoryParams(ctx, req)
	return s.commands.Create(ctx, params)
}

func (s *CategoryService) GetByID(ctx context.Context, id string) (*models.Category, error) {
	ctx, span := observability.StartServiceSpan(ctx, "CategoryService", "GetByID")
	defer span.End()
	return s.query.GetByID(ctx, id)
}

func (s *CategoryService) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	ctx, span := observability.StartServiceSpan(ctx, "CategoryService", "GetIDBySubID")
	defer span.End()
	return s.query.GetIDBySubID(ctx, subID)
}

func (s *CategoryService) List(ctx context.Context, params *models.CategorySearchParams) ([]*models.Category, int64, error) {
	ctx, span := observability.StartServiceSpan(ctx, "CategoryService", "List")
	defer span.End()

	queryParams := valueobjects.ToListCategoriesParams(ctx, params)
	return s.query.List(ctx, queryParams)
}

func (s *CategoryService) Update(ctx context.Context, id string, req *models.UpdateCategoryRequest) (*models.Category, error) {
	ctx, span := observability.StartServiceSpan(ctx, "CategoryService", "Update")
	defer span.End()

	params := valueobjects.ToUpdateCategoryParams(ctx, req)
	return s.commands.Update(ctx, id, params)
}

func (s *CategoryService) Delete(ctx context.Context, id string) (*models.Category, error) {
	ctx, span := observability.StartServiceSpan(ctx, "CategoryService", "Delete")
	defer span.End()
	return s.commands.Delete(ctx, id)
}
