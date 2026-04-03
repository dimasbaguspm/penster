package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type RateCurrencyRepository struct {
	db *query.Queries
}

func NewRateCurrencyRepository(db *query.Queries) *RateCurrencyRepository {
	return &RateCurrencyRepository{db: db}
}

func (r *RateCurrencyRepository) Upsert(ctx context.Context, req *models.UpsertRateCurrencyRequest) (*models.RateCurrency, error) {
	result, err := r.db.UpsertRateCurrency(ctx, query.UpsertRateCurrencyParams{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         req.Rate,
		RateDate:     req.RateDate,
	})
	if err != nil {
		return nil, err
	}
	return toRateCurrencyModel(result), nil
}

func (r *RateCurrencyRepository) Get(ctx context.Context, from, to string, date time.Time) (*models.RateCurrency, error) {
	result, err := r.db.GetRateCurrency(ctx, query.GetRateCurrencyParams{
		FromCurrency: from,
		ToCurrency:   to,
		RateDate:     date,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toRateCurrencyModel(result), nil
}

func (r *RateCurrencyRepository) List(ctx context.Context, params *models.RateCurrencySearchParams) ([]*models.RateCurrency, int64, error) {
	rows, err := r.db.ListRateCurrencies(ctx, query.ListRateCurrenciesParams{
		FromCurrency: conv.StringPtrToEmpty(params.FromCurrency),
		ToCurrency:   conv.StringPtrToEmpty(params.ToCurrency),
		PageSize:     int64(params.PageSize),
		Offset:       int64(params.Offset()),
	})
	if err != nil {
		return nil, 0, err
	}

	currencies := make([]*models.RateCurrency, 0, len(rows))
	var total int64
	for _, row := range rows {
		currencies = append(currencies, toRateCurrencyModel(query.RateCurrency{
			ID:           row.ID,
			FromCurrency: row.FromCurrency,
			ToCurrency:   row.ToCurrency,
			Rate:         row.Rate,
			RateDate:     row.RateDate,
			CreatedAt:    row.CreatedAt,
		}))
		total = row.Total
	}

	return currencies, total, nil
}

func (r *RateCurrencyRepository) Prune(ctx context.Context, olderThan time.Time) (int64, error) {
	return r.db.PruneOldRates(ctx, olderThan)
}

func toRateCurrencyModel(q query.RateCurrency) *models.RateCurrency {
	m := &models.RateCurrency{
		FromCurrency: q.FromCurrency,
		ToCurrency:   q.ToCurrency,
		RateDate:     time.Time{},
		CreatedAt:    time.Time{},
	}

	if q.ID.Valid {
		m.ID = uuid.UUID(q.ID.Bytes).String()
	}

	if q.Rate.Valid {
		rate, _ := q.Rate.Float64Value()
		m.Rate = rate.Float64
	}

	if q.RateDate.Valid {
		m.RateDate = q.RateDate.Time
	}

	if q.CreatedAt.Valid {
		m.CreatedAt = q.CreatedAt.Time
	}

	return m
}
