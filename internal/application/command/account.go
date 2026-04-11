package command

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

// AccountCommandInterface defines write operations for accounts
type AccountCommandInterface interface {
	Create(ctx context.Context, params query.CreateAccountParams) (*models.Account, error)
	Update(ctx context.Context, id string, params query.UpdateAccountParams) (*models.Account, error)
	Delete(ctx context.Context, id string) (*models.Account, error)
	UpdateBalance(ctx context.Context, id string, newBalance int64) (*models.Account, error)
}

var _ AccountCommandInterface = (*AccountCommand)(nil)

type AccountCommand struct {
	repo *repository.AccountRepository
}

// NewAccountCommand creates a new AccountCommand
func NewAccountCommand(repo *repository.AccountRepository) *AccountCommand {
	return &AccountCommand{repo: repo}
}

func (c *AccountCommand) Create(ctx context.Context, params query.CreateAccountParams) (*models.Account, error) {
	ctx, span := observability.StartCommandSpan(ctx, "account", "create")
	defer span.End()
	return c.repo.Create(ctx, params)
}

func (c *AccountCommand) Update(ctx context.Context, id string, params query.UpdateAccountParams) (*models.Account, error) {
	ctx, span := observability.StartCommandSpan(ctx, "account", "update")
	defer span.End()
	return c.repo.UpdateBySubID(ctx, id, params)
}

func (c *AccountCommand) Delete(ctx context.Context, id string) (*models.Account, error) {
	ctx, span := observability.StartCommandSpan(ctx, "account", "delete")
	defer span.End()
	return c.repo.DeleteBySubID(ctx, id)
}

func (c *AccountCommand) UpdateBalance(ctx context.Context, id string, newBalance int64) (*models.Account, error) {
	ctx, span := observability.StartCommandSpan(ctx, "account", "update_balance")
	defer span.End()
	internalID, err := c.repo.GetIDBySubID(ctx, id)
	if err != nil {
		return nil, err
	}
	if internalID == 0 {
		return nil, nil
	}
	return c.repo.UpdateBalanceByID(ctx, internalID, newBalance)
}
