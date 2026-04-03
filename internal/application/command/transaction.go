package command

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type TransactionCommandInterface interface {
	Create(ctx context.Context, req *models.CreateTransactionRequest, currencyRate float64) (*models.Transaction, error)
	Update(ctx context.Context, id string, req *models.UpdateTransactionRequest, currencyRate float64) (*models.Transaction, error)
	Delete(ctx context.Context, id string) (*models.Transaction, error)
}

var _ TransactionCommandInterface = (*TransactionCommand)(nil)

type TransactionCommand struct {
	repo *repository.TransactionRepository
}

func NewTransactionCommand(repo *repository.TransactionRepository) *TransactionCommand {
	return &TransactionCommand{repo: repo}
}

func (c *TransactionCommand) Create(ctx context.Context, req *models.CreateTransactionRequest, currencyRate float64) (*models.Transaction, error) {
	return c.repo.Create(ctx, req, currencyRate)
}

func (c *TransactionCommand) Update(ctx context.Context, id string, req *models.UpdateTransactionRequest, currencyRate float64) (*models.Transaction, error) {
	return c.repo.UpdateBySubID(ctx, id, req, currencyRate)
}

func (c *TransactionCommand) Delete(ctx context.Context, id string) (*models.Transaction, error) {
	return c.repo.DeleteBySubID(ctx, id)
}
