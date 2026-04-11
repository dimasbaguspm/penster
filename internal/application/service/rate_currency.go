package service

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/application/command"
	"github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/domain/valueobjects"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type RateCurrencyService struct {
	query    query.RateCurrencyQueryInterface
	commands command.RateCurrencyCommandInterface
}

func NewRateCurrencyService(query query.RateCurrencyQueryInterface, commands command.RateCurrencyCommandInterface) *RateCurrencyService {
	return &RateCurrencyService{
		query:    query,
		commands: commands,
	}
}

func (s *RateCurrencyService) Upsert(ctx context.Context, req *models.UpsertRateCurrencyRequest) (*models.RateCurrency, error) {
	ctx, span := observability.StartServiceSpan(ctx, "RateCurrencyService", "Upsert")
	defer span.End()

	params := valueobjects.ToUpsertRateCurrencyParams(ctx, req)
	return s.commands.Upsert(ctx, params)
}

func (s *RateCurrencyService) Get(ctx context.Context, from, to string, rateDate time.Time) (*models.RateCurrency, error) {
	ctx, span := observability.StartServiceSpan(ctx, "RateCurrencyService", "Get")
	defer span.End()
	return s.query.Get(ctx, from, to, rateDate)
}

func (s *RateCurrencyService) GetRate(ctx context.Context, currency, baseCurrency string) (float64, error) {
	ctx, span := observability.StartServiceSpan(ctx, "RateCurrencyService", "GetRate")
	defer span.End()
	return s.query.GetRate(ctx, currency, baseCurrency)
}

func (s *RateCurrencyService) List(ctx context.Context, params *models.RateCurrencySearchParams) ([]*models.RateCurrency, int64, error) {
	ctx, span := observability.StartServiceSpan(ctx, "RateCurrencyService", "List")
	defer span.End()

	queryParams := valueobjects.ToListRateCurrenciesParams(ctx, params)
	return s.query.List(ctx, queryParams)
}

func (s *RateCurrencyService) Prune(ctx context.Context, olderThan time.Time) (int64, error) {
	ctx, span := observability.StartServiceSpan(ctx, "RateCurrencyService", "Prune")
	defer span.End()
	return s.commands.Prune(ctx, olderThan)
}
