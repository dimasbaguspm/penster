package query

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type RateCurrencyQueryInterface interface {
	Get(ctx context.Context, from, to string, rateDate time.Time) (*models.RateCurrency, error)
	List(ctx context.Context, params *models.RateCurrencySearchParams) ([]*models.RateCurrency, int64, error)
}

var _ RateCurrencyQueryInterface = (*RateCurrencyQuery)(nil)

type RateCurrencyQuery struct {
	repo *repository.RateCurrencyRepository
}

func NewRateCurrencyQuery(repo *repository.RateCurrencyRepository) *RateCurrencyQuery {
	return &RateCurrencyQuery{repo: repo}
}

func (q *RateCurrencyQuery) Get(ctx context.Context, from, to string, rateDate time.Time) (*models.RateCurrency, error) {
	return q.repo.Get(ctx, from, to, rateDate)
}

func (q *RateCurrencyQuery) List(ctx context.Context, params *models.RateCurrencySearchParams) ([]*models.RateCurrency, int64, error) {
	return q.repo.List(ctx, params)
}
