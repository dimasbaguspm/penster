package command

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type TransactionCommandInterface interface {
	Create(ctx context.Context, params query.CreateTransactionParams) (*models.Transaction, error)
	Update(ctx context.Context, id string, params query.UpdateTransactionParams) (*models.Transaction, error)
	Delete(ctx context.Context, id string) (*models.Transaction, error)
}

var _ TransactionCommandInterface = (*TransactionCommand)(nil)

type TransactionCommand struct {
	repo *repository.TransactionRepository
}

func NewTransactionCommand(repo *repository.TransactionRepository) *TransactionCommand {
	return &TransactionCommand{repo: repo}
}

func (c *TransactionCommand) Create(ctx context.Context, params query.CreateTransactionParams) (*models.Transaction, error) {
	ctx, span := observability.StartCommandSpan(ctx, "transaction", "create")
	defer span.End()
	return c.repo.Create(ctx, params)
}

func (c *TransactionCommand) Update(ctx context.Context, id string, params query.UpdateTransactionParams) (*models.Transaction, error) {
	ctx, span := observability.StartCommandSpan(ctx, "transaction", "update")
	defer span.End()
	return c.repo.UpdateBySubID(ctx, id, params)
}

func (c *TransactionCommand) Delete(ctx context.Context, id string) (*models.Transaction, error) {
	ctx, span := observability.StartCommandSpan(ctx, "transaction", "delete")
	defer span.End()
	return c.repo.DeleteBySubID(ctx, id)
}
