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
	log := observability.NewLogger(ctx, "service", "rate_currency")
	ctx, span := observability.StartServiceSpan(log.Context(), "rate_currency", "Upsert")
	defer span.End()

	log.Info("upsert started", "from", req.FromCurrency, "to", req.ToCurrency)
	result, err := s.commands.Upsert(ctx, valueobjects.ToUpsertRateCurrencyParams(ctx, req))
	if err != nil {
		log.Error("upsert failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("upsert succeeded", "from", req.FromCurrency, "to", req.ToCurrency)
	return result, nil
}

func (s *RateCurrencyService) Get(ctx context.Context, from, to string, rateDate time.Time) (*models.RateCurrency, error) {
	log := observability.NewLogger(ctx, "service", "rate_currency")
	ctx, span := observability.StartServiceSpan(log.Context(), "rate_currency", "Get")
	defer span.End()

	log.Info("get started", "from", from, "to", to, "date", rateDate.Format("2006-01-02"))
	return s.query.Get(ctx, from, to, rateDate)
}

func (s *RateCurrencyService) GetRate(ctx context.Context, currency, baseCurrency string) (float64, error) {
	log := observability.NewLogger(ctx, "service", "rate_currency")
	ctx, span := observability.StartServiceSpan(log.Context(), "rate_currency", "GetRate")
	defer span.End()

	log.Info("get_rate started", "currency", currency, "base_currency", baseCurrency)
	return s.query.GetRate(ctx, currency, baseCurrency)
}

func (s *RateCurrencyService) List(ctx context.Context, params *models.RateCurrencySearchParams) ([]*models.RateCurrency, int64, error) {
	log := observability.NewLogger(ctx, "service", "rate_currency")
	ctx, span := observability.StartServiceSpan(log.Context(), "rate_currency", "List")
	defer span.End()

	log.Info("list started")
	queryParams := valueobjects.ToListRateCurrenciesParams(ctx, params)
	currencies, total, err := s.query.List(ctx, queryParams)
	if err != nil {
		log.Error("list failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, 0, err
	}
	log.Info("list succeeded", "count", len(currencies), "total", total)
	return currencies, total, nil
}

func (s *RateCurrencyService) Prune(ctx context.Context, olderThan time.Time) (int64, error) {
	log := observability.NewLogger(ctx, "service", "rate_currency")
	ctx, span := observability.StartServiceSpan(log.Context(), "rate_currency", "Prune")
	defer span.End()

	log.Info("prune started", "older_than", olderThan.Format("2006-01-02"))
	result, err := s.commands.Prune(ctx, olderThan)
	if err != nil {
		log.Error("prune failed", "error", err)
		observability.RecordError(ctx, err)
		return 0, err
	}
	log.Info("prune succeeded", "deleted", result)
	return result, nil
}
