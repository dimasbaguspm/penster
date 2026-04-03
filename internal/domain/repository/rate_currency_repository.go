package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type RateCurrencyRepository struct {
	db *query.Queries
}

func NewRateCurrencyRepository(db *query.Queries) *RateCurrencyRepository {
	return &RateCurrencyRepository{db: db}
}

func (r *RateCurrencyRepository) Upsert(ctx context.Context, req *models.UpsertRateCurrencyRequest) (*models.RateCurrency, error) {
	var rateNumeric pgtype.Numeric
	if err := rateNumeric.Scan(fmt.Sprintf("%.6f", req.Rate)); err != nil {
		return nil, err
	}

	result, err := r.db.UpsertRateCurrency(ctx, query.UpsertRateCurrencyParams{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rateNumeric,
		RateDate:     pgtype.Date{Time: req.RateDate, Valid: true},
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
		RateDate:     pgtype.Date{Time: date, Valid: true},
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
	var fromCurrency, toCurrency string
	if params.FromCurrency != nil {
		fromCurrency = *params.FromCurrency
	}
	if params.ToCurrency != nil {
		toCurrency = *params.ToCurrency
	}

	rows, err := r.db.ListRateCurrencies(ctx, query.ListRateCurrenciesParams{
		Column1: fromCurrency,
		Column2: toCurrency,
		Column3: params.PageSize,
		Offset:  int32(params.Offset()),
	})
	if err != nil {
		return nil, 0, err
	}

	currencies := make([]*models.RateCurrency, 0, len(rows))
	for _, row := range rows {
		currencies = append(currencies, toRateCurrencyModel(row))
	}

	total, err := r.db.CountRateCurrencies(ctx, query.CountRateCurrenciesParams{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	})
	if err != nil {
		return nil, 0, err
	}

	return currencies, total, nil
}

func (r *RateCurrencyRepository) Prune(ctx context.Context, olderThan time.Time) error {
	return r.db.PruneOldRates(ctx, pgtype.Date{Time: olderThan, Valid: true})
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
