package query

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type TransactionQueryInterface interface {
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	List(ctx context.Context, params *models.TransactionSearchParams) ([]*models.Transaction, int64, error)
}

var _ TransactionQueryInterface = (*TransactionQuery)(nil)

type TransactionQuery struct {
	repo *repository.TransactionRepository
}

func NewTransactionQuery(repo *repository.TransactionRepository) *TransactionQuery {
	return &TransactionQuery{repo: repo}
}

func (q *TransactionQuery) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	return q.repo.GetBySubID(ctx, id)
}

func (q *TransactionQuery) List(ctx context.Context, params *models.TransactionSearchParams) ([]*models.Transaction, int64, error) {
	return q.repo.List(ctx, params)
}
