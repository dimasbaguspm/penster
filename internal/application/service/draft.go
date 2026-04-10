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

type draftIDs struct {
	accountID         int32
	transferAccountID pgtype.Int4
	categoryID        int32
}

func (s *DraftService) validateRelatedEntities(ctx context.Context, accountID, transferAccountID, categoryID string) (*draftIDs, error) {
	var ids draftIDs

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

type DraftService struct {
	query               appquery.DraftQueryInterface
	commands            command.DraftCommandInterface
	accountService      *AccountService
	categoryService     *CategoryService
	rateCurrencyService *RateCurrencyService
	transactionService  *TransactionService
	cfg                 *config.Config
}

func NewDraftService(
	query appquery.DraftQueryInterface,
	commands command.DraftCommandInterface,
	accountService *AccountService,
	categoryService *CategoryService,
	rateCurrencyService *RateCurrencyService,
	transactionService *TransactionService,
	cfg *config.Config,
) *DraftService {
	return &DraftService{
		query:               query,
		commands:            commands,
		accountService:      accountService,
		categoryService:     categoryService,
		rateCurrencyService: rateCurrencyService,
		transactionService:  transactionService,
		cfg:                 cfg,
	}
}

func (s *DraftService) Create(ctx context.Context, req *models.CreateDraftRequest) (*models.Draft, error) {
	ctx, span := observability.StartServiceSpan(ctx, "DraftService", "Create")
	defer span.End()

	if req.TransactionType == string(models.TransactionTypeTransfer) && req.AccountID == req.TransferAccountID {
		return nil, entities.ErrTransferAccountNotFound
	}

	ids, err := s.validateRelatedEntities(ctx, req.AccountID, req.TransferAccountID, req.CategoryID)
	if err != nil {
		return nil, err
	}

	currencyRate, err := s.rateCurrencyService.GetRate(ctx, req.Currency, s.cfg.App.BaseCurrency)
	if err != nil {
		return nil, err
	}

	params, err := valueobjects.ToCreateDraftParams(ids.accountID, ids.transferAccountID, ids.categoryID, currencyRate, req)
	if err != nil {
		return nil, err
	}

	return s.commands.Create(ctx, params)
}

func (s *DraftService) GetByID(ctx context.Context, id string) (*models.Draft, error) {
	ctx, span := observability.StartServiceSpan(ctx, "DraftService", "GetByID")
	defer span.End()
	return s.query.GetByID(ctx, id)
}

func (s *DraftService) List(ctx context.Context, params *models.DraftSearchParams) ([]*models.Draft, int64, error) {
	ctx, span := observability.StartServiceSpan(ctx, "DraftService", "List")
	defer span.End()
	queryParams := valueobjects.ToListDraftsParams(params)
	return s.query.List(ctx, queryParams)
}

func (s *DraftService) Update(ctx context.Context, id string, req *models.UpdateDraftRequest) (*models.Draft, error) {
	ctx, span := observability.StartServiceSpan(ctx, "DraftService", "Update")
	defer span.End()

	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, entities.ErrDraftNotFound
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

	// Validate same-account transfer
	transactionType := existing.TransactionType
	if req.TransactionType != nil {
		transactionType = *req.TransactionType
	}
	if transactionType == string(models.TransactionTypeTransfer) && accountID == transferAccountID {
		return nil, entities.ErrTransferAccountNotFound
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

	updateParams := valueobjects.ToUpdateDraftParams(ids.accountID, ids.transferAccountID, ids.categoryID, currencyRate, req)

	return s.commands.Update(ctx, id, updateParams)
}

func (s *DraftService) Confirm(ctx context.Context, draftSubID string) (*models.Transaction, error) {
	ctx, span := observability.StartServiceSpan(ctx, "DraftService", "Confirm")
	defer span.End()

	draft, err := s.GetByID(ctx, draftSubID)
	if err != nil {
		return nil, err
	}
	if draft == nil {
		return nil, entities.ErrDraftNotFound
	}
	if draft.Status != string(models.DraftStatusPending) {
		return nil, entities.ErrDraftNotPending
	}

	txReq := &models.CreateTransactionRequest{
		AccountID:         draft.AccountID,
		TransferAccountID: draft.TransferAccountID,
		CategoryID:        draft.CategoryID,
		TransactionType:   models.TransactionType(draft.TransactionType),
		Title:             draft.Title,
		Amount:            draft.Amount,
		Currency:          draft.Currency,
		Notes:             draft.Notes,
	}

	tx, err := s.transactionService.Create(ctx, txReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction from draft: %w", err)
	}

	// Update draft status to confirmed
	if err := s.commands.UpdateStatus(ctx, draftSubID, string(models.DraftStatusConfirmed)); err != nil {
		return nil, fmt.Errorf("failed to update draft status: %w", err)
	}

	return tx, nil
}

func (s *DraftService) Reject(ctx context.Context, draftSubID string) error {
	ctx, span := observability.StartServiceSpan(ctx, "DraftService", "Reject")
	defer span.End()

	draft, err := s.GetByID(ctx, draftSubID)
	if err != nil {
		return err
	}
	if draft == nil {
		return entities.ErrDraftNotFound
	}
	if draft.Status != string(models.DraftStatusPending) {
		return entities.ErrDraftNotPending
	}

	return s.commands.UpdateStatus(ctx, draftSubID, string(models.DraftStatusRejected))
}

func (s *DraftService) Delete(ctx context.Context, draftSubID string) error {
	ctx, span := observability.StartServiceSpan(ctx, "DraftService", "Delete")
	defer span.End()

	draft, err := s.GetByID(ctx, draftSubID)
	if err != nil {
		return err
	}
	if draft == nil {
		return entities.ErrDraftNotFound
	}
	if draft.Status != string(models.DraftStatusRejected) {
		return entities.ErrDraftNotRejected
	}

	return s.commands.Delete(ctx, draftSubID)
}
