package query

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type RateCurrencyQueryInterface interface {
	Get(ctx context.Context, from, to string, rateDate time.Time) (*models.RateCurrency, error)
	GetRate(ctx context.Context, currency, baseCurrency string) (float64, error)
	List(ctx context.Context, params query.ListRateCurrenciesParams) ([]*models.RateCurrency, int64, error)
}

var _ RateCurrencyQueryInterface = (*RateCurrencyQuery)(nil)

type RateCurrencyQuery struct {
	repo *repository.RateCurrencyRepository
}

func NewRateCurrencyQuery(repo *repository.RateCurrencyRepository) *RateCurrencyQuery {
	return &RateCurrencyQuery{repo: repo}
}

func (q *RateCurrencyQuery) Get(ctx context.Context, from, to string, rateDate time.Time) (*models.RateCurrency, error) {
	ctx, span := observability.StartQuerySpan(ctx, "rate_currency", "get")
	defer span.End()
	return q.repo.Get(ctx, from, to, rateDate)
}

func (q *RateCurrencyQuery) GetRate(ctx context.Context, currency, baseCurrency string) (float64, error) {
	ctx, span := observability.StartQuerySpan(ctx, "rate_currency", "get_rate")
	defer span.End()
	if currency == baseCurrency {
		return 1, nil
	}
	rateDate := time.Now().Truncate(24 * time.Hour)
	rate, err := q.Get(ctx, currency, baseCurrency, rateDate)
	if err != nil {
		return 0, err
	}
	if rate == nil {
		return 1, nil
	}
	return rate.Rate, nil
}

func (q *RateCurrencyQuery) List(ctx context.Context, params query.ListRateCurrenciesParams) ([]*models.RateCurrency, int64, error) {
	ctx, span := observability.StartQuerySpan(ctx, "rate_currency", "list")
	defer span.End()
	return q.repo.List(ctx, params)
}
