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
	log := observability.NewLogger(ctx, "service", "transaction")
	ctx, span := observability.StartServiceSpan(log.Context(), "transaction", "Create")
	defer span.End()

	log.Info("create started", "account_id", req.AccountID, "category_id", req.CategoryID, "amount", req.Amount)

	log.Debug("create validating_related_entities", "account_id", req.AccountID, "transfer_account_id", req.TransferAccountID, "category_id", req.CategoryID)
	ids, err := s.validateRelatedEntities(ctx, req.AccountID, req.TransferAccountID, req.CategoryID)
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Debug("create entities_validated", "account_id", ids.accountID)

	// Validate transfer doesn't result in negative balance
	if req.TransactionType == models.TransactionTypeTransfer {
		log.Debug("create validating_transfer", "account_id", req.AccountID, "amount", req.Amount)
		if err := s.accountService.ValidateTransfer(ctx, req.AccountID, req.Amount); err != nil {
			log.Error("create failed", "error", err)
			observability.RecordError(ctx, err)
			return nil, err
		}
		log.Debug("create transfer_validated")
	}

	log.Debug("create fetching_currency_rate", "currency", req.Currency, "base_currency", s.cfg.App.BaseCurrency)
	currencyRate, err := s.rateCurrencyService.GetRate(ctx, req.Currency, s.cfg.App.BaseCurrency)
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Debug("create currency_rate_fetched", "rate", currencyRate)

	params, err := valueobjects.ToCreateTransactionParams(ids.accountID, ids.transferAccountID, ids.categoryID, currencyRate, req)
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	tx, err := s.commands.Create(ctx, params)
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Debug("create transaction_saved", "tx_id", tx.ID)

	log.Debug("create updating_balances", "account_id", req.AccountID, "type", tx.TransactionType, "amount", tx.Amount)
	if err := s.accountService.UpdateAccountBalances(ctx, req.AccountID, req.TransferAccountID, tx.TransactionType, tx.Amount); err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}
	log.Debug("create balances_updated")

	log.Info("create succeeded", "id", tx.ID)
	observability.TransactionsCreated.Add(ctx, 1)
	return tx, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "service", "transaction")
	ctx, span := observability.StartServiceSpan(log.Context(), "transaction", "GetByID")
	defer span.End()

	log.Info("get_by_id started", "id", id)
	return s.query.GetByID(ctx, id)
}

func (s *TransactionService) List(ctx context.Context, params *models.TransactionSearchParams) ([]*models.Transaction, int64, error) {
	log := observability.NewLogger(ctx, "service", "transaction")
	ctx, span := observability.StartServiceSpan(log.Context(), "transaction", "List")
	defer span.End()

	log.Info("list started")
	queryParams := valueobjects.ToListTransactionsParams(params)
	transactions, total, err := s.query.List(ctx, queryParams)
	if err != nil {
		log.Error("list failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, 0, err
	}
	log.Info("list succeeded", "count", len(transactions), "total", total)
	return transactions, total, nil
}

func (s *TransactionService) Update(ctx context.Context, id string, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "service", "transaction")
	ctx, span := observability.StartServiceSpan(log.Context(), "transaction", "Update")
	defer span.End()

	log.Info("update started", "id", id)

	existing, err := s.GetByID(ctx, id)
	if err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	log.Debug("update fetching_existing", "id", id, "type", existing.TransactionType, "amount", existing.Amount)

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
			err := entities.ErrTransferToSameAccount
			log.Error("update failed", "error", err)
			observability.RecordError(ctx, err)
			return nil, err
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

	log.Debug("update validating_entities", "account_id", accountID, "transfer_account_id", transferAccountID, "category_id", categoryID)
	ids, err := s.validateRelatedEntities(ctx, accountID, transferAccountID, categoryID)
	if err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Debug("update entities_validated")

	currencyRate := existing.CurrencyRate
	if req.Currency != nil && *req.Currency != existing.Currency {
		log.Debug("update fetching_currency_rate", "currency", *req.Currency, "base_currency", s.cfg.App.BaseCurrency)
		currencyRate, err = s.rateCurrencyService.GetRate(ctx, *req.Currency, s.cfg.App.BaseCurrency)
		if err != nil {
			log.Error("update failed", "error", err)
			observability.RecordError(ctx, err)
			return nil, err
		}
		log.Debug("update currency_rate_fetched", "rate", currencyRate)
	}

	log.Debug("update reverse_balances", "account_id", existing.AccountID, "type", existing.TransactionType, "amount", existing.Amount)
	if err := s.accountService.ReverseAccountBalances(ctx, existing.AccountID, existing.TransferAccountID, existing.TransactionType, existing.Amount); err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
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
			log.Error("update failed", "error", err)
			observability.RecordError(ctx, err)
			return nil, err
		}
		if accountAfterReversal.Balance < newAmount {
			err := entities.ErrInsufficientBalance
			log.Error("update failed", "error", err)
			observability.RecordError(ctx, err)
			return nil, err
		}
	}

	updateParams := valueobjects.ToUpdateTransactionParams(ids.accountID, ids.transferAccountID, ids.categoryID, currencyRate, req)

	tx, err := s.commands.Update(ctx, id, updateParams)
	if err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
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

	log.Debug("update apply_balances", "account_id", accountIDStr, "type", tx.TransactionType, "amount", tx.Amount)
	if err := s.accountService.UpdateAccountBalances(ctx, accountIDStr, transferAccountIDStr, tx.TransactionType, tx.Amount); err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	log.Info("update succeeded", "id", id)
	observability.TransactionsUpdated.Add(ctx, 1)
	return tx, nil
}

func (s *TransactionService) Delete(ctx context.Context, id string) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "service", "transaction")
	ctx, span := observability.StartServiceSpan(log.Context(), "transaction", "Delete")
	defer span.End()

	log.Info("delete started", "id", id)

	existing, err := s.GetByID(ctx, id)
	if err != nil {
		log.Error("delete failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	log.Debug("delete reverse_balances", "account_id", existing.AccountID, "type", existing.TransactionType, "amount", existing.Amount)
	if err := s.accountService.ReverseAccountBalances(ctx, existing.AccountID, existing.TransferAccountID, existing.TransactionType, existing.Amount); err != nil {
		log.Error("delete failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, fmt.Errorf("failed to reverse account balance: %w", err)
	}

	result, err := s.commands.Delete(ctx, id)
	if err != nil {
		log.Error("delete failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("delete succeeded", "id", id)
	observability.TransactionsDeleted.Add(ctx, 1)
	return result, nil
}
