package service

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/command"
	"github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type TransactionService struct {
	query               query.TransactionQueryInterface
	commands            command.TransactionCommandInterface
	rateCurrencyService *RateCurrencyService
	cfg                 *config.Config
}

func NewTransactionService(
	query query.TransactionQueryInterface,
	commands command.TransactionCommandInterface,
	rateCurrencyService *RateCurrencyService,
	cfg *config.Config,
) *TransactionService {
	return &TransactionService{
		query:               query,
		commands:            commands,
		rateCurrencyService: rateCurrencyService,
		cfg:                 cfg,
	}
}

func (s *TransactionService) Create(ctx context.Context, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	currencyRate := float64(1)
	if req.Currency != s.cfg.App.BaseCurrency {
		rateDate := time.Now().Truncate(24 * time.Hour)
		rate, err := s.rateCurrencyService.Get(ctx, req.Currency, s.cfg.App.BaseCurrency, rateDate)
		if err != nil {
			return nil, err
		}
		if rate != nil {
			currencyRate = rate.Rate
		}
	}

	return s.commands.Create(ctx, req, currencyRate)
}

func (s *TransactionService) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	return s.query.GetByID(ctx, id)
}

func (s *TransactionService) List(ctx context.Context, params *models.TransactionSearchParams) ([]*models.Transaction, int64, error) {
	return s.query.List(ctx, params)
}

func (s *TransactionService) Update(ctx context.Context, id string, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	currencyRate := existing.CurrencyRate
	if req.Currency != nil && *req.Currency != existing.Currency {
		rateDate := time.Now().Truncate(24 * time.Hour)
		rate, err := s.rateCurrencyService.Get(ctx, *req.Currency, s.cfg.App.BaseCurrency, rateDate)
		if err != nil {
			return nil, err
		}
		if rate != nil {
			currencyRate = rate.Rate
		}
	}

	return s.commands.Update(ctx, id, req, currencyRate)
}

func (s *TransactionService) Delete(ctx context.Context, id string) (*models.Transaction, error) {
	return s.commands.Delete(ctx, id)
}
