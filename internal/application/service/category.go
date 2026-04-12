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
	log := observability.NewLogger(ctx, "service", "category")
	ctx, span := observability.StartServiceSpan(log.Context(), "category", "Create")
	defer span.End()

	log.Info("create started", "name", req.Name)
	result, err := s.commands.Create(ctx, valueobjects.ToCreateCategoryParams(ctx, req))
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("create succeeded", "id", result.ID)
	observability.CategoriesCreated.Add(ctx, 1)
	return result, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id string) (*models.Category, error) {
	log := observability.NewLogger(ctx, "service", "category")
	ctx, span := observability.StartServiceSpan(log.Context(), "category", "GetByID")
	defer span.End()

	log.Info("get_by_id started", "id", id)
	return s.query.GetByID(ctx, id)
}

func (s *CategoryService) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	ctx, span := observability.StartServiceSpan(ctx, "CategoryService", "GetIDBySubID")
	defer span.End()
	return s.query.GetIDBySubID(ctx, subID)
}

func (s *CategoryService) List(ctx context.Context, params *models.CategorySearchParams) ([]*models.Category, int64, error) {
	log := observability.NewLogger(ctx, "service", "category")
	ctx, span := observability.StartServiceSpan(log.Context(), "category", "List")
	defer span.End()

	log.Info("list started")
	queryParams := valueobjects.ToListCategoriesParams(ctx, params)
	categories, total, err := s.query.List(ctx, queryParams)
	if err != nil {
		log.Error("list failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, 0, err
	}
	log.Info("list succeeded", "count", len(categories), "total", total)
	return categories, total, nil
}

func (s *CategoryService) Update(ctx context.Context, id string, req *models.UpdateCategoryRequest) (*models.Category, error) {
	log := observability.NewLogger(ctx, "service", "category")
	ctx, span := observability.StartServiceSpan(log.Context(), "category", "Update")
	defer span.End()

	log.Info("update started", "id", id)
	params := valueobjects.ToUpdateCategoryParams(ctx, req)
	result, err := s.commands.Update(ctx, id, params)
	if err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("update succeeded", "id", id)
	observability.CategoriesUpdated.Add(ctx, 1)
	return result, nil
}

func (s *CategoryService) Delete(ctx context.Context, id string) (*models.Category, error) {
	log := observability.NewLogger(ctx, "service", "category")
	ctx, span := observability.StartServiceSpan(log.Context(), "category", "Delete")
	defer span.End()

	log.Info("delete started", "id", id)
	result, err := s.commands.Delete(ctx, id)
	if err != nil {
		log.Error("delete failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("delete succeeded", "id", id)
	observability.CategoriesDeleted.Add(ctx, 1)
	return result, nil
}
