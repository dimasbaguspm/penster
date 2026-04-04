package query

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type RateCurrencyQueryInterface interface {
	Get(ctx context.Context, from, to string, rateDate time.Time) (*models.RateCurrency, error)
	GetRate(ctx context.Context, currency, baseCurrency string) (float64, error)
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

func (q *RateCurrencyQuery) GetRate(ctx context.Context, currency, baseCurrency string) (float64, error) {
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

func (q *RateCurrencyQuery) List(ctx context.Context, params *models.RateCurrencySearchParams) ([]*models.RateCurrency, int64, error) {
	return q.repo.List(ctx, params)
}
