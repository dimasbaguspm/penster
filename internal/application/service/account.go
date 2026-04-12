package service

import (
	"context"
	"fmt"

	"github.com/dimasbaguspm/penster/internal/application/command"
	"github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/internal/domain/valueobjects"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type AccountService struct {
	query    query.AccountQueryInterface
	commands command.AccountCommandInterface
}

func NewAccountService(query query.AccountQueryInterface, commands command.AccountCommandInterface) *AccountService {
	return &AccountService{
		query:    query,
		commands: commands,
	}
}

func (s *AccountService) Create(ctx context.Context, req *models.CreateAccountRequest) (*models.Account, error) {
	log := observability.NewLogger(ctx, "service", "account")
	ctx, span := observability.StartServiceSpan(log.Context(), "account", "Create")
	defer span.End()

	log.Info("create started", "name", req.Name)
	result, err := s.commands.Create(ctx, valueobjects.ToCreateAccountParams(ctx, req))
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("create succeeded", "id", result.ID)
	return result, nil
}

func (s *AccountService) GetByID(ctx context.Context, id string) (*models.Account, error) {
	log := observability.NewLogger(ctx, "service", "account")
	ctx, span := observability.StartServiceSpan(log.Context(), "account", "GetByID")
	defer span.End()

	log.Info("get_by_id started", "id", id)
	return s.query.GetByID(ctx, id)
}

func (s *AccountService) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	ctx, span := observability.StartServiceSpan(ctx, "AccountService", "GetIDBySubID")
	defer span.End()
	return s.query.GetIDBySubID(ctx, subID)
}

func (s *AccountService) List(ctx context.Context, params *models.AccountSearchParams) ([]*models.Account, int64, error) {
	log := observability.NewLogger(ctx, "service", "account")
	ctx, span := observability.StartServiceSpan(log.Context(), "account", "List")
	defer span.End()

	log.Info("list started")
	queryParams := valueobjects.ToListAccountsParams(ctx, params)
	accounts, total, err := s.query.List(ctx, queryParams)
	if err != nil {
		log.Error("list failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, 0, err
	}
	log.Info("list succeeded", "count", len(accounts), "total", total)
	return accounts, total, nil
}

func (s *AccountService) Update(ctx context.Context, id string, req *models.UpdateAccountRequest) (*models.Account, error) {
	log := observability.NewLogger(ctx, "service", "account")
	ctx, span := observability.StartServiceSpan(log.Context(), "account", "Update")
	defer span.End()

	log.Info("update started", "id", id)
	params := valueobjects.ToUpdateAccountParams(ctx, req)
	result, err := s.commands.Update(ctx, id, params)
	if err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("update succeeded", "id", id)
	return result, nil
}

func (s *AccountService) Delete(ctx context.Context, id string) (*models.Account, error) {
	log := observability.NewLogger(ctx, "service", "account")
	ctx, span := observability.StartServiceSpan(log.Context(), "account", "Delete")
	defer span.End()

	log.Info("delete started", "id", id)
	result, err := s.commands.Delete(ctx, id)
	if err != nil {
		log.Error("delete failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("delete succeeded", "id", id)
	return result, nil
}

func (s *AccountService) UpdateAccountBalances(ctx context.Context, accountID string, transferAccountID string, transactionType models.TransactionType, amount int64) error {
	log := observability.NewLogger(ctx, "service", "account")
	ctx, span := observability.StartServiceSpan(log.Context(), "account", "UpdateAccountBalances")
	defer span.End()

	log.Info("update_balances started", "account_id", accountID, "type", transactionType, "amount", amount)

	account, err := s.GetByID(ctx, accountID)
	if err != nil {
		log.Error("update_balances failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	if account == nil {
		err := fmt.Errorf("account not found: %s", accountID)
		log.Error("update_balances failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}

	var newBalance int64
	switch transactionType {
	case models.TransactionTypeExpense, models.TransactionTypeTransfer:
		newBalance = account.Balance - amount
	case models.TransactionTypeIncome:
		newBalance = account.Balance + amount
	}

	log.Debug("update_balances calculating", "current_balance", account.Balance, "type", transactionType, "amount", amount, "new_balance", newBalance)

	_, err = s.commands.UpdateBalance(ctx, accountID, newBalance)
	if err != nil {
		log.Error("update_balances failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	log.Debug("update_balances source_updated", "account_id", accountID, "new_balance", newBalance)

	if transactionType == models.TransactionTypeTransfer && transferAccountID != "" && transferAccountID != accountID {
		log.Debug("update_balances fetching_destination", "transfer_account_id", transferAccountID)
		destAccount, err := s.GetByID(ctx, transferAccountID)
		if err != nil {
			log.Error("update_balances failed", "error", err)
			observability.RecordError(ctx, err)
			return err
		}
		if destAccount == nil {
			err := fmt.Errorf("transfer account not found: %s", transferAccountID)
			log.Error("update_balances failed", "error", err)
			observability.RecordError(ctx, err)
			return err
		}

		newDestBalance := destAccount.Balance + amount
		log.Debug("update_balances updating_destination", "transfer_account_id", transferAccountID, "current_balance", destAccount.Balance, "new_balance", newDestBalance)
		_, err = s.commands.UpdateBalance(ctx, transferAccountID, newDestBalance)
		if err != nil {
			log.Error("update_balances failed", "error", err)
			observability.RecordError(ctx, err)
			return err
		}
		log.Debug("update_balances destination_updated")
	}

	log.Info("update_balances succeeded", "account_id", accountID, "new_balance", newBalance)
	return nil
}

func (s *AccountService) ReverseAccountBalances(ctx context.Context, accountID string, transferAccountID string, transactionType models.TransactionType, amount int64) error {
	log := observability.NewLogger(ctx, "service", "account")
	ctx, span := observability.StartServiceSpan(log.Context(), "account", "ReverseAccountBalances")
	defer span.End()

	log.Info("reverse_balances started", "account_id", accountID, "type", transactionType, "amount", amount)

	account, err := s.GetByID(ctx, accountID)
	if err != nil {
		log.Error("reverse_balances failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	if account == nil {
		err := fmt.Errorf("account not found: %s", accountID)
		log.Error("reverse_balances failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}

	var newBalance int64
	switch transactionType {
	case models.TransactionTypeExpense, models.TransactionTypeTransfer:
		newBalance = account.Balance + amount
	case models.TransactionTypeIncome:
		newBalance = account.Balance - amount
	}

	log.Debug("reverse_balances calculating", "current_balance", account.Balance, "type", transactionType, "amount", amount, "new_balance", newBalance)

	_, err = s.commands.UpdateBalance(ctx, accountID, newBalance)
	if err != nil {
		log.Error("reverse_balances failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	log.Debug("reverse_balances source_updated", "account_id", accountID, "new_balance", newBalance)

	if transactionType == models.TransactionTypeTransfer && transferAccountID != "" && transferAccountID != accountID {
		log.Debug("reverse_balances fetching_destination", "transfer_account_id", transferAccountID)
		destAccount, err := s.GetByID(ctx, transferAccountID)
		if err != nil {
			log.Error("reverse_balances failed", "error", err)
			observability.RecordError(ctx, err)
			return err
		}
		if destAccount == nil {
			err := fmt.Errorf("transfer account not found: %s", transferAccountID)
			log.Error("reverse_balances failed", "error", err)
			observability.RecordError(ctx, err)
			return err
		}

		newDestBalance := destAccount.Balance - amount
		log.Debug("reverse_balances updating_destination", "transfer_account_id", transferAccountID, "current_balance", destAccount.Balance, "new_balance", newDestBalance)
		_, err = s.commands.UpdateBalance(ctx, transferAccountID, newDestBalance)
		if err != nil {
			log.Error("reverse_balances failed", "error", err)
			observability.RecordError(ctx, err)
			return err
		}
		log.Debug("reverse_balances destination_updated")
	}

	log.Info("reverse_balances succeeded", "account_id", accountID, "new_balance", newBalance)
	return nil
}

// ValidateTransfer checks if a transfer would result in negative balance on the source account.
func (s *AccountService) ValidateTransfer(ctx context.Context, accountID string, amount int64) error {
	log := observability.NewLogger(ctx, "service", "account")
	ctx, span := observability.StartServiceSpan(log.Context(), "account", "ValidateTransfer")
	defer span.End()

	log.Info("validate_transfer started", "account_id", accountID, "amount", amount)

	account, err := s.GetByID(ctx, accountID)
	if err != nil {
		log.Error("validate_transfer failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	if account == nil {
		err := fmt.Errorf("account not found: %s", accountID)
		log.Error("validate_transfer failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}

	if account.Balance < amount {
		log.Debug("validate_transfer insufficient_balance", "balance", account.Balance, "amount", amount)
		return entities.ErrInsufficientBalance
	}

	log.Info("validate_transfer succeeded", "account_id", accountID, "balance", account.Balance, "amount", amount)
	return nil
}
