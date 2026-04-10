package service

import (
	"context"
	"fmt"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/command"
	appquery "github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/internal/domain/valueobjects"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
	"github.com/dimasbaguspm/penster/pkg/syncerr"
	"github.com/jackc/pgx/v5/pgtype"
)

type transactionIDs struct {
	accountID         int32
	transferAccountID pgtype.Int4
	categoryID        int32
}

func (s *TransactionService) validateRelatedEntities(ctx context.Context, accountID, transferAccountID, categoryID string) (*transactionIDs, error) {
	var ids transactionIDs

	grp := syncerr.Group{}

	grp.Go(func() error {
		id, err := s.accountService.GetIDBySubID(ctx, accountID)
		if err != nil {
			return err
		}
		if id == 0 {
			return fmt.Errorf("%w: %s", entities.ErrAccountNotFound, accountID)
		}
		ids.accountID = id
		return nil
	})

	if transferAccountID != "" {
		grp.Go(func() error {
			id, err := s.accountService.GetIDBySubID(ctx, transferAccountID)
			if err != nil {
				return err
			}
			if id == 0 {
				return fmt.Errorf("%w: %s", entities.ErrTransferAccountNotFound, transferAccountID)
			}
			ids.transferAccountID = pgtype.Int4{Int32: id, Valid: true}
			return nil
		})
	}

	grp.Go(func() error {
		id, err := s.categoryService.GetIDBySubID(ctx, categoryID)
		if err != nil {
			return err
		}
		if id == 0 {
			return fmt.Errorf("%w: %s", entities.ErrCategoryNotFound, categoryID)
		}
		ids.categoryID = id
		return nil
	})

	if errs := grp.Wait(); len(errs) > 0 {
		return nil, errs[0]
	}

	return &ids, nil
}

type TransactionService struct {
	query               appquery.TransactionQueryInterface
	commands            command.TransactionCommandInterface
	accountService      *AccountService
	categoryService     *CategoryService
	rateCurrencyService *RateCurrencyService
	cfg                 *config.Config
}

func NewTransactionService(
	query appquery.TransactionQueryInterface,
	commands command.TransactionCommandInterface,
	accountService *AccountService,
	categoryService *CategoryService,
	rateCurrencyService *RateCurrencyService,
	cfg *config.Config,
) *TransactionService {
	return &TransactionService{
		query:               query,
		commands:            commands,
		accountService:      accountService,
		categoryService:     categoryService,
		rateCurrencyService: rateCurrencyService,
		cfg:                 cfg,
	}
}

func (s *TransactionService) Create(ctx context.Context, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	ctx, span := observability.StartServiceSpan(ctx, "TransactionService", "Create")
	defer span.End()

	ids, err := s.validateRelatedEntities(ctx, req.AccountID, req.TransferAccountID, req.CategoryID)
	if err != nil {
		return nil, err
	}

	// Validate transfer doesn't result in negative balance
	if req.TransactionType == models.TransactionTypeTransfer {
		if err := s.accountService.ValidateTransfer(ctx, req.AccountID, req.Amount); err != nil {
			return nil, err
		}
	}

	currencyRate, err := s.rateCurrencyService.GetRate(ctx, req.Currency, s.cfg.App.BaseCurrency)
	if err != nil {
		return nil, err
	}

	params, err := valueobjects.ToCreateTransactionParams(ids.accountID, ids.transferAccountID, ids.categoryID, currencyRate, req)
	if err != nil {
		return nil, err
	}

	tx, err := s.commands.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	if err := s.accountService.UpdateAccountBalances(ctx, req.AccountID, req.TransferAccountID, tx.TransactionType, tx.Amount); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	return tx, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	ctx, span := observability.StartServiceSpan(ctx, "TransactionService", "GetByID")
	defer span.End()
	return s.query.GetByID(ctx, id)
}

func (s *TransactionService) List(ctx context.Context, params *models.TransactionSearchParams) ([]*models.Transaction, int64, error) {
	ctx, span := observability.StartServiceSpan(ctx, "TransactionService", "List")
	defer span.End()
	queryParams := valueobjects.ToListTransactionsParams(params)
	return s.query.List(ctx, queryParams)
}

func (s *TransactionService) Update(ctx context.Context, id string, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	ctx, span := observability.StartServiceSpan(ctx, "TransactionService", "Update")
	defer span.End()

	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

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
			return nil, entities.ErrTransferToSameAccount
		}
	}

	accountID := existing.AccountID
	transferAccountID := existing.TransferAccountID
	categoryID := existing.CategoryID

	if req.AccountID != nil {
		accountID = *req.AccountID
	}
	if req.TransferAccountID != nil {
		transferAccountID = *req.TransferAccountID
	}
	if req.CategoryID != nil {
		categoryID = *req.CategoryID
	}

	ids, err := s.validateRelatedEntities(ctx, accountID, transferAccountID, categoryID)
	if err != nil {
		return nil, err
	}

	currencyRate := existing.CurrencyRate
	if req.Currency != nil && *req.Currency != existing.Currency {
		currencyRate, err = s.rateCurrencyService.GetRate(ctx, *req.Currency, s.cfg.App.BaseCurrency)
		if err != nil {
			return nil, err
		}
	}

	if err := s.accountService.ReverseAccountBalances(ctx, existing.AccountID, existing.TransferAccountID, existing.TransactionType, existing.Amount); err != nil {
		return nil, fmt.Errorf("failed to reverse account balance: %w", err)
	}

	// Determine the new transaction type for validation
	newTransactionType := existing.TransactionType
	if req.TransactionType != nil {
		newTransactionType = *req.TransactionType
	}

	// Validate transfer doesn't result in negative balance after reversal
	newAmount := existing.Amount
	if req.Amount != nil {
		newAmount = *req.Amount
	}
	if newTransactionType == models.TransactionTypeTransfer {
		accountAfterReversal, err := s.accountService.GetByID(ctx, accountID)
		if err != nil {
			return nil, err
		}
		if accountAfterReversal.Balance < newAmount {
			return nil, entities.ErrInsufficientBalance
		}
	}

	updateParams := valueobjects.ToUpdateTransactionParams(ids.accountID, ids.transferAccountID, ids.categoryID, currencyRate, req)

	tx, err := s.commands.Update(ctx, id, updateParams)
	if err != nil {
		return nil, err
	}

	accountIDStr := existing.AccountID
	if req.AccountID != nil {
		accountIDStr = *req.AccountID
	}
	var transferAccountIDStr string
	if req.TransferAccountID != nil {
		transferAccountIDStr = *req.TransferAccountID
	} else {
		transferAccountIDStr = existing.TransferAccountID
	}

	if err := s.accountService.UpdateAccountBalances(ctx, accountIDStr, transferAccountIDStr, tx.TransactionType, tx.Amount); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	return tx, nil
}

func (s *TransactionService) Delete(ctx context.Context, id string) (*models.Transaction, error) {
	ctx, span := observability.StartServiceSpan(ctx, "TransactionService", "Delete")
	defer span.End()

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
