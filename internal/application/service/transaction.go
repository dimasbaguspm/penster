package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/command"
	"github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type TransactionService struct {
	query               query.TransactionQueryInterface
	commands            command.TransactionCommandInterface
	accountService      *AccountService
	rateCurrencyService *RateCurrencyService
	cfg                 *config.Config
}

func NewTransactionService(
	query query.TransactionQueryInterface,
	commands command.TransactionCommandInterface,
	accountService *AccountService,
	rateCurrencyService *RateCurrencyService,
	cfg *config.Config,
) *TransactionService {
	return &TransactionService{
		query:               query,
		commands:            commands,
		accountService:      accountService,
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

	tx, err := s.commands.Create(ctx, req, currencyRate)
	if err != nil {
		return nil, err
	}

	if err := s.accountService.UpdateAccountBalances(ctx, req.AccountID, req.TransferAccountID, tx.TransactionType, tx.Amount); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	return tx, nil
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

	// Check for same-account transfer
	// Check if updating TransactionType to transfer, or if already a transfer and updating transfer_account_id
	isTransfer := existing.TransactionType == models.TransactionTypeTransfer
	if req.TransactionType != nil {
		isTransfer = *req.TransactionType == models.TransactionTypeTransfer
	}
	if isTransfer && req.TransferAccountID != nil {
		accountID := existing.AccountID
		if req.AccountID != nil {
			accountID = *req.AccountID
		}
		if accountID == *req.TransferAccountID {
			return nil, repository.ErrTransferToSameAccount
		}
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

	if err := s.accountService.ReverseAccountBalances(ctx, existing.AccountID, existing.TransferAccountID, existing.TransactionType, existing.Amount); err != nil {
		return nil, fmt.Errorf("failed to reverse account balance: %w", err)
	}

	tx, err := s.commands.Update(ctx, id, req, currencyRate)
	if err != nil {
		return nil, err
	}

	accountID := existing.AccountID
	if req.AccountID != nil {
		accountID = *req.AccountID
	}
	var transferAccountID *string
	if req.TransferAccountID != nil {
		transferAccountID = req.TransferAccountID
	} else {
		transferAccountID = &existing.TransferAccountID
	}

	if err := s.accountService.UpdateAccountBalances(ctx, accountID, *transferAccountID, tx.TransactionType, tx.Amount); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	return tx, nil
}

func (s *TransactionService) Delete(ctx context.Context, id string) (*models.Transaction, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	if err := s.accountService.ReverseAccountBalances(ctx, existing.AccountID, existing.TransferAccountID, existing.TransactionType, existing.Amount); err != nil {
		return nil, fmt.Errorf("failed to reverse account balance: %w", err)
	}

	return s.commands.Delete(ctx, id)
}
