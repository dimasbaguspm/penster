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
	ctx, span := observability.StartServiceSpan(ctx, "DraftService", "validateRelatedEntities")
	defer span.End()

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
	log := observability.NewLogger(ctx, "service", "draft")
	ctx, span := observability.StartServiceSpan(log.Context(), "draft", "Create")
	defer span.End()

	log.Info("create started", "title", req.Title, "type", req.TransactionType)

	if req.TransactionType == string(models.TransactionTypeTransfer) && req.AccountID == req.TransferAccountID {
		err := entities.ErrTransferAccountNotFound
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	ids, err := s.validateRelatedEntities(ctx, req.AccountID, req.TransferAccountID, req.CategoryID)
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	currencyRate, err := s.rateCurrencyService.GetRate(ctx, req.Currency, s.cfg.App.BaseCurrency)
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	params, err := valueobjects.ToCreateDraftParams(ctx, ids.accountID, ids.transferAccountID, ids.categoryID, currencyRate, req)
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	result, err := s.commands.Create(ctx, params)
	if err != nil {
		log.Error("create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("create succeeded", "id", result.ID)
	return result, nil
}

func (s *DraftService) GetByID(ctx context.Context, id string) (*models.Draft, error) {
	log := observability.NewLogger(ctx, "service", "draft")
	ctx, span := observability.StartServiceSpan(log.Context(), "draft", "GetByID")
	defer span.End()

	log.Info("get_by_id started", "id", id)
	return s.query.GetByID(ctx, id)
}

func (s *DraftService) List(ctx context.Context, params *models.DraftSearchParams) ([]*models.Draft, int64, error) {
	log := observability.NewLogger(ctx, "service", "draft")
	ctx, span := observability.StartServiceSpan(log.Context(), "draft", "List")
	defer span.End()

	log.Info("list started")
	queryParams := valueobjects.ToListDraftsParams(ctx, params)
	drafts, total, err := s.query.List(ctx, queryParams)
	if err != nil {
		log.Error("list failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, 0, err
	}
	log.Info("list succeeded", "count", len(drafts), "total", total)
	return drafts, total, nil
}

func (s *DraftService) Update(ctx context.Context, id string, req *models.UpdateDraftRequest) (*models.Draft, error) {
	log := observability.NewLogger(ctx, "service", "draft")
	ctx, span := observability.StartServiceSpan(log.Context(), "draft", "Update")
	defer span.End()

	log.Info("update started", "id", id)

	existing, err := s.GetByID(ctx, id)
	if err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
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
		err := entities.ErrTransferAccountNotFound
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	ids, err := s.validateRelatedEntities(ctx, accountID, transferAccountID, categoryID)
	if err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	currencyRate := existing.CurrencyRate
	if req.Currency != nil && *req.Currency != existing.Currency {
		currencyRate, err = s.rateCurrencyService.GetRate(ctx, *req.Currency, s.cfg.App.BaseCurrency)
		if err != nil {
			log.Error("update failed", "error", err)
			observability.RecordError(ctx, err)
			return nil, err
		}
	}

	updateParams := valueobjects.ToUpdateDraftParams(ctx, ids.accountID, ids.transferAccountID, ids.categoryID, currencyRate, req)

	result, err := s.commands.Update(ctx, id, updateParams)
	if err != nil {
		log.Error("update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("update succeeded", "id", id)
	return result, nil
}

func (s *DraftService) Confirm(ctx context.Context, draftSubID string) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "service", "draft")
	ctx, span := observability.StartServiceSpan(log.Context(), "draft", "Confirm")
	defer span.End()

	log.Info("confirm started", "draft_sub_id", draftSubID)

	draft, err := s.GetByID(ctx, draftSubID)
	if err != nil {
		log.Error("confirm failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	if draft == nil {
		err := entities.ErrDraftNotFound
		log.Error("confirm failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Debug("confirm draft_fetched", "status", draft.Status, "type", draft.TransactionType)

	if draft.Status != string(models.DraftStatusPending) {
		err := entities.ErrDraftNotPending
		log.Error("confirm failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Debug("confirm status_validated")

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

	log.Debug("confirm creating_transaction", "account_id", draft.AccountID, "category_id", draft.CategoryID, "amount", draft.Amount)
	tx, err := s.transactionService.Create(ctx, txReq)
	if err != nil {
		log.Error("confirm failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, fmt.Errorf("failed to create transaction from draft: %w", err)
	}
	log.Debug("confirm transaction_created", "tx_id", tx.ID)

	log.Debug("confirm updating_status", "draft_sub_id", draftSubID, "status", "confirmed")
	if err := s.commands.UpdateStatus(ctx, draftSubID, string(models.DraftStatusConfirmed)); err != nil {
		log.Error("confirm failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, fmt.Errorf("failed to update draft status: %w", err)
	}
	log.Debug("confirm status_updated")

	log.Info("confirm succeeded", "draft_sub_id", draftSubID, "tx_id", tx.ID)
	return tx, nil
}

func (s *DraftService) Reject(ctx context.Context, draftSubID string) error {
	log := observability.NewLogger(ctx, "service", "draft")
	ctx, span := observability.StartServiceSpan(log.Context(), "draft", "Reject")
	defer span.End()

	log.Info("reject started", "draft_sub_id", draftSubID)

	draft, err := s.GetByID(ctx, draftSubID)
	if err != nil {
		log.Error("reject failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	if draft == nil {
		err := entities.ErrDraftNotFound
		log.Error("reject failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	log.Debug("reject draft_fetched", "status", draft.Status)

	if draft.Status != string(models.DraftStatusPending) {
		err := entities.ErrDraftNotPending
		log.Error("reject failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	log.Debug("reject status_validated")

	log.Debug("reject updating_status", "draft_sub_id", draftSubID, "status", "rejected")
	err = s.commands.UpdateStatus(ctx, draftSubID, string(models.DraftStatusRejected))
	if err != nil {
		log.Error("reject failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	log.Debug("reject status_updated")
	log.Info("reject succeeded", "draft_sub_id", draftSubID)
	return nil
}

func (s *DraftService) Delete(ctx context.Context, draftSubID string) error {
	log := observability.NewLogger(ctx, "service", "draft")
	ctx, span := observability.StartServiceSpan(log.Context(), "draft", "Delete")
	defer span.End()

	log.Info("delete started", "draft_sub_id", draftSubID)

	draft, err := s.GetByID(ctx, draftSubID)
	if err != nil {
		log.Error("delete failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	if draft == nil {
		err := entities.ErrDraftNotFound
		log.Error("delete failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	if draft.Status != string(models.DraftStatusRejected) {
		err := entities.ErrDraftNotRejected
		log.Error("delete failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}

	err = s.commands.Delete(ctx, draftSubID)
	if err != nil {
		log.Error("delete failed", "error", err)
		observability.RecordError(ctx, err)
		return err
	}
	log.Info("delete succeeded", "draft_sub_id", draftSubID)
	return nil
}
