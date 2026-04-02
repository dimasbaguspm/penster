package service

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/application/command"
	"github.com/dimasbaguspm/penster/internal/application/query"
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

func (s *AccountService) List(ctx context.Context, params *models.AccountSearchParams) ([]*models.Account, int64, error) {
	return s.query.List(ctx, params)
}

func (s *AccountService) Update(ctx context.Context, id string, req *models.UpdateAccountRequest) (*models.Account, error) {
	return s.commands.Update(ctx, id, req)
}

func (s *AccountService) Delete(ctx context.Context, id string) (*models.Account, error) {
	return s.commands.Delete(ctx, id)
}
