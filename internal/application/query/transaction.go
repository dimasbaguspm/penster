package query

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type TransactionQueryInterface interface {
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	List(ctx context.Context, params query.ListTransactionsParams) ([]*models.Transaction, int64, error)
}

var _ TransactionQueryInterface = (*TransactionQuery)(nil)

type TransactionQuery struct {
	repo *repository.TransactionRepository
}

func NewTransactionQuery(repo *repository.TransactionRepository) *TransactionQuery {
	return &TransactionQuery{repo: repo}
}

func (q *TransactionQuery) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	ctx, span := observability.StartQuerySpan(ctx, "transaction", "get_by_id")
	defer span.End()
	return q.repo.GetBySubID(ctx, id)
}

func (q *TransactionQuery) List(ctx context.Context, params query.ListTransactionsParams) ([]*models.Transaction, int64, error) {
	ctx, span := observability.StartQuerySpan(ctx, "transaction", "list")
	defer span.End()
	return q.repo.List(ctx, params)
}
