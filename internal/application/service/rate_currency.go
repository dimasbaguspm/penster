package service

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/application/command"
	"github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/pkg/models"
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
	return s.commands.Upsert(ctx, req)
}

func (s *RateCurrencyService) Get(ctx context.Context, from, to string, rateDate time.Time) (*models.RateCurrency, error) {
	return s.query.Get(ctx, from, to, rateDate)
}

func (s *RateCurrencyService) List(ctx context.Context, params *models.RateCurrencySearchParams) ([]*models.RateCurrency, int64, error) {
	return s.query.List(ctx, params)
}

func (s *RateCurrencyService) Prune(ctx context.Context, olderThan time.Time) (int64, error) {
	return s.commands.Prune(ctx, olderThan)
}
