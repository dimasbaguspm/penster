package query

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

// AccountQueryInterface defines read operations for accounts
type AccountQueryInterface interface {
	GetByID(ctx context.Context, id string) (*models.Account, error)
	GetIDBySubID(ctx context.Context, subID string) (int32, error)
	List(ctx context.Context, params query.ListAccountsParams) ([]*models.Account, int64, error)
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
	ctx, span := observability.StartQuerySpan(ctx, "account", "get_by_id")
	defer span.End()
	return q.repo.GetBySubID(ctx, id)
}

func (q *AccountQuery) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	ctx, span := observability.StartQuerySpan(ctx, "account", "get_id_by_sub_id")
	defer span.End()
	return q.repo.GetIDBySubID(ctx, subID)
}

func (q *AccountQuery) List(ctx context.Context, params query.ListAccountsParams) ([]*models.Account, int64, error) {
	ctx, span := observability.StartQuerySpan(ctx, "account", "list")
	defer span.End()
	return q.repo.List(ctx, params)
}
