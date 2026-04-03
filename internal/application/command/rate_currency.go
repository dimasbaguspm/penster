package command

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type RateCurrencyCommandInterface interface {
	Upsert(ctx context.Context, req *models.UpsertRateCurrencyRequest) (*models.RateCurrency, error)
	Prune(ctx context.Context, olderThan time.Time) (int64, error)
}

var _ RateCurrencyCommandInterface = (*RateCurrencyCommand)(nil)

type RateCurrencyCommand struct {
	repo *repository.RateCurrencyRepository
}

func NewRateCurrencyCommand(repo *repository.RateCurrencyRepository) *RateCurrencyCommand {
	return &RateCurrencyCommand{repo: repo}
}

func (c *RateCurrencyCommand) Upsert(ctx context.Context, req *models.UpsertRateCurrencyRequest) (*models.RateCurrency, error) {
	return c.repo.Upsert(ctx, req)
}

func (c *RateCurrencyCommand) Prune(ctx context.Context, olderThan time.Time) (int64, error) {
	return c.repo.Prune(ctx, olderThan)
}
