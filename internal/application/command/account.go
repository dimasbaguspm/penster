package command

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

// AccountCommandInterface defines write operations for accounts
type AccountCommandInterface interface {
	Create(ctx context.Context, req *models.CreateAccountRequest) (*models.Account, error)
	Update(ctx context.Context, id string, req *models.UpdateAccountRequest) (*models.Account, error)
	Delete(ctx context.Context, id string) (*models.Account, error)
}

var _ AccountCommandInterface = (*AccountCommand)(nil)

type AccountCommand struct {
	repo *repository.AccountRepository
}

// NewAccountCommand creates a new AccountCommand
func NewAccountCommand(repo *repository.AccountRepository) *AccountCommand {
	return &AccountCommand{repo: repo}
}

func (c *AccountCommand) Create(ctx context.Context, req *models.CreateAccountRequest) (*models.Account, error) {
	return c.repo.Create(ctx, req)
}

func (c *AccountCommand) Update(ctx context.Context, id string, req *models.UpdateAccountRequest) (*models.Account, error) {
	return c.repo.UpdateBySubID(ctx, id, req)
}

func (c *AccountCommand) Delete(ctx context.Context, id string) (*models.Account, error) {
	return c.repo.DeleteBySubID(ctx, id)
}
