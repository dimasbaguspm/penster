package command

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type RateCurrencyCommandInterface interface {
	Upsert(ctx context.Context, params query.UpsertRateCurrencyParams) (*models.RateCurrency, error)
	Prune(ctx context.Context, olderThan time.Time) (int64, error)
}

var _ RateCurrencyCommandInterface = (*RateCurrencyCommand)(nil)

type RateCurrencyCommand struct {
	repo *repository.RateCurrencyRepository
}

func NewRateCurrencyCommand(repo *repository.RateCurrencyRepository) *RateCurrencyCommand {
	return &RateCurrencyCommand{repo: repo}
}

func (c *RateCurrencyCommand) Upsert(ctx context.Context, params query.UpsertRateCurrencyParams) (*models.RateCurrency, error) {
	ctx, span := observability.StartCommandSpan(ctx, "rate_currency", "upsert")
	defer span.End()
	return c.repo.Upsert(ctx, params)
}

func (c *RateCurrencyCommand) Prune(ctx context.Context, olderThan time.Time) (int64, error) {
	ctx, span := observability.StartCommandSpan(ctx, "rate_currency", "prune")
	defer span.End()
	err := c.repo.Prune(ctx, olderThan)
	if err != nil {
		return 0, err
	}
	return 0, nil
}
