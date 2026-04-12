package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type RateCurrencyRepository struct {
	db *query.Queries
}

func NewRateCurrencyRepository(db *query.Queries) *RateCurrencyRepository {
	return &RateCurrencyRepository{db: db}
}

func (r *RateCurrencyRepository) Upsert(ctx context.Context, params query.UpsertRateCurrencyParams) (*models.RateCurrency, error) {
	log := observability.NewLogger(ctx, "repository", "rate_currency")
	ctx, span := observability.StartRepoSpan(log.Context(), "rate_currency", "Upsert")
	defer span.End()

	log.Info("rate_currency.upsert started", "from", params.FromCurrency, "to", params.ToCurrency)
	result, err := r.db.UpsertRateCurrency(ctx, params)
	if err != nil {
		log.Error("rate_currency.upsert failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("rate_currency.upsert succeeded")
	return toRateCurrencyModel(ctx, result), nil
}

func (r *RateCurrencyRepository) Get(ctx context.Context, from, to string, date time.Time) (*models.RateCurrency, error) {
	log := observability.NewLogger(ctx, "repository", "rate_currency")
	ctx, span := observability.StartRepoSpan(log.Context(), "rate_currency", "Get")
	defer span.End()

	log.Info("rate_currency.get started", "from", from, "to", to, "date", date.Format("2006-01-02"))
	result, err := r.db.GetRateCurrency(ctx, query.GetRateCurrencyParams{
		FromCurrency: from,
		ToCurrency:   to,
		RateDate:     pgtype.Date{Time: date, Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("rate_currency.get not found", "from", from, "to", to)
			return nil, nil
		}
		log.Error("rate_currency.get failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("rate_currency.get succeeded")
	return toRateCurrencyModel(ctx, result), nil
}

func (r *RateCurrencyRepository) List(ctx context.Context, params query.ListRateCurrenciesParams) ([]*models.RateCurrency, int64, error) {
	log := observability.NewLogger(ctx, "repository", "rate_currency")
	ctx, span := observability.StartRepoSpan(log.Context(), "rate_currency", "List")
	defer span.End()

	log.Info("rate_currency.list started")
	rows, err := r.db.ListRateCurrencies(ctx, params)
	if err != nil {
		log.Error("rate_currency.list failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, 0, err
	}

	currencies := make([]*models.RateCurrency, 0, len(rows))
	for _, row := range rows {
		currencies = append(currencies, toRateCurrencyModel(ctx, row))
	}

	total, err := r.db.CountRateCurrencies(ctx, query.CountRateCurrenciesParams{
		FromCurrency: params.Column1,
		ToCurrency:   params.Column2,
	})
	if err != nil {
		log.Error("rate_currency.list.count failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, 0, err
	}

	log.Info("rate_currency.list succeeded", "count", len(currencies), "total", total)
	return currencies, total, nil
}

func (r *RateCurrencyRepository) Prune(ctx context.Context, olderThan time.Time) error {
	log := observability.NewLogger(ctx, "repository", "rate_currency")
	ctx, span := observability.StartRepoSpan(log.Context(), "rate_currency", "Prune")
	defer span.End()

	log.Info("rate_currency.prune started", "older_than", olderThan.Format("2006-01-02"))
	err := r.db.PruneOldRates(ctx, pgtype.Date{Time: olderThan, Valid: true})
	if err != nil {
		log.Error("rate_currency.prune failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	log.Info("rate_currency.prune succeeded")
	return nil
}

func toRateCurrencyModel(ctx context.Context, q query.RateCurrency) *models.RateCurrency {
	_, span := observability.StartRepoSpan(ctx, "rate_currencies", "to_model")
	defer span.End()

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
