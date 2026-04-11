package query

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type DraftQueryInterface interface {
	GetByID(ctx context.Context, id string) (*models.Draft, error)
	List(ctx context.Context, params query.ListDraftsParams) ([]*models.Draft, int64, error)
}

var _ DraftQueryInterface = (*DraftQuery)(nil)

type DraftQuery struct {
	repo *repository.DraftRepository
}

func NewDraftQuery(repo *repository.DraftRepository) *DraftQuery {
	return &DraftQuery{repo: repo}
}

func (q *DraftQuery) GetByID(ctx context.Context, id string) (*models.Draft, error) {
	ctx, span := observability.StartQuerySpan(ctx, "draft", "get_by_id")
	defer span.End()
	return q.repo.GetBySubID(ctx, id)
}

func (q *DraftQuery) List(ctx context.Context, params query.ListDraftsParams) ([]*models.Draft, int64, error) {
	ctx, span := observability.StartQuerySpan(ctx, "draft", "list")
	defer span.End()
	return q.repo.List(ctx, params)
}
