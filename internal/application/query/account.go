package query

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

// AccountQueryInterface defines read operations for accounts
type AccountQueryInterface interface {
	GetByID(ctx context.Context, id string) (*models.Account, error)
	GetIDBySubID(ctx context.Context, subID string) (int32, error)
	List(ctx context.Context, params *models.AccountSearchParams) ([]*models.Account, int64, error)
}

var _ AccountQueryInterface = (*AccountQuery)(nil)

// AccountQuery implements AccountQueryInterface using repository
type AccountQuery struct {
	repo *repository.AccountRepository
}

// NewAccountQuery creates a new AccountQuery
func NewAccountQuery(repo *repository.AccountRepository) *AccountQuery {
	return &AccountQuery{repo: repo}
}

func (q *AccountQuery) GetByID(ctx context.Context, id string) (*models.Account, error) {
	return q.repo.GetBySubID(ctx, id)
}

func (q *AccountQuery) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	return q.repo.GetIDBySubID(ctx, subID)
}

func (q *AccountQuery) List(ctx context.Context, params *models.AccountSearchParams) ([]*models.Account, int64, error) {
	return q.repo.List(ctx, params)
}
