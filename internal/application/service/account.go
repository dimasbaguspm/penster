package service

import (
	"context"
	"fmt"

	"github.com/dimasbaguspm/penster/internal/application/command"
	"github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/pkg/models"
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
	return s.commands.Create(ctx, req)
}

func (s *AccountService) GetByID(ctx context.Context, id string) (*models.Account, error) {
	return s.query.GetByID(ctx, id)
}

func (s *AccountService) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	return s.query.GetIDBySubID(ctx, subID)
}

func (s *AccountService) List(ctx context.Context, params *models.AccountSearchParams) ([]*models.Account, int64, error) {
	return s.query.List(ctx, params)
}

func (s *AccountService) Update(ctx context.Context, id string, req *models.UpdateAccountRequest) (*models.Account, error) {
	return s.commands.Update(ctx, id, req)
}

func (s *AccountService) Delete(ctx context.Context, id string) (*models.Account, error) {
	return s.commands.Delete(ctx, id)
}

func (s *AccountService) UpdateAccountBalances(ctx context.Context, accountID string, transferAccountID string, transactionType models.TransactionType, amount int64) error {
	account, err := s.GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("account not found: %s", accountID)
	}

	var newBalance int64
	switch transactionType {
	case models.TransactionTypeExpense, models.TransactionTypeTransfer:
		newBalance = account.Balance - amount
	case models.TransactionTypeIncome:
		newBalance = account.Balance + amount
	}

	_, err = s.commands.UpdateBalance(ctx, accountID, newBalance)
	if err != nil {
		return err
	}

	if transactionType == models.TransactionTypeTransfer && transferAccountID != "" && transferAccountID != accountID {
		destAccount, err := s.GetByID(ctx, transferAccountID)
		if err != nil {
			return err
		}
		if destAccount == nil {
			return fmt.Errorf("transfer account not found: %s", transferAccountID)
		}

		newDestBalance := destAccount.Balance + amount
		_, err = s.commands.UpdateBalance(ctx, transferAccountID, newDestBalance)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AccountService) ReverseAccountBalances(ctx context.Context, accountID string, transferAccountID string, transactionType models.TransactionType, amount int64) error {
	account, err := s.GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("account not found: %s", accountID)
	}

	var newBalance int64
	switch transactionType {
	case models.TransactionTypeExpense, models.TransactionTypeTransfer:
		newBalance = account.Balance + amount
	case models.TransactionTypeIncome:
		newBalance = account.Balance - amount
	}

	_, err = s.commands.UpdateBalance(ctx, accountID, newBalance)
	if err != nil {
		return err
	}

	if transactionType == models.TransactionTypeTransfer && transferAccountID != "" && transferAccountID != accountID {
		destAccount, err := s.GetByID(ctx, transferAccountID)
		if err != nil {
			return err
		}
		if destAccount == nil {
			return fmt.Errorf("transfer account not found: %s", transferAccountID)
		}

		newDestBalance := destAccount.Balance - amount
		_, err = s.commands.UpdateBalance(ctx, transferAccountID, newDestBalance)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateTransfer checks if a transfer would result in negative balance on the source account.
func (s *AccountService) ValidateTransfer(ctx context.Context, accountID string, amount int64) error {
	account, err := s.GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("account not found: %s", accountID)
	}

	if account.Balance < amount {
		return entities.ErrInsufficientBalance
	}

	return nil
}
