package command

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type DraftCommandInterface interface {
	Create(ctx context.Context, params query.CreateDraftParams) (*models.Draft, error)
	Update(ctx context.Context, id string, params query.UpdateDraftParams) (*models.Draft, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
}

var _ DraftCommandInterface = (*DraftCommand)(nil)

type DraftCommand struct {
	repo *repository.DraftRepository
}

func NewDraftCommand(repo *repository.DraftRepository) *DraftCommand {
	return &DraftCommand{repo: repo}
}

func (c *DraftCommand) Create(ctx context.Context, params query.CreateDraftParams) (*models.Draft, error) {
	return c.repo.Create(ctx, params)
}

func (c *DraftCommand) Update(ctx context.Context, id string, params query.UpdateDraftParams) (*models.Draft, error) {
	return c.repo.UpdateBySubID(ctx, id, params)
}

func (c *DraftCommand) UpdateStatus(ctx context.Context, id string, status string) error {
	return c.repo.UpdateStatus(ctx, id, status)
}

func (c *DraftCommand) Delete(ctx context.Context, id string) error {
	return c.repo.SoftDelete(ctx, id)
}
