package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type DraftRepository struct {
	db *query.Queries
}

func NewDraftRepository(db *query.Queries) *DraftRepository {
	return &DraftRepository{db: db}
}

func (r *DraftRepository) Create(ctx context.Context, params query.CreateDraftParams) (*models.Draft, error) {
	log := observability.NewLogger(ctx, "repository", "draft")
	ctx, span := observability.StartRepoSpan(log.Context(), "draft", "Create")
	defer span.End()

	log.Info("draft.create started", "title", params.Title)
	id, err := r.db.CreateDraft(ctx, params)
	if err != nil {
		log.Error("draft.create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("draft.created", "id", id)
	return r.GetByID(ctx, id)
}

func (r *DraftRepository) GetByID(ctx context.Context, id int32) (*models.Draft, error) {
	log := observability.NewLogger(ctx, "repository", "draft")
	ctx, span := observability.StartRepoSpan(log.Context(), "draft", "GetByID")
	defer span.End()

	log.Info("draft.get_by_id started", "id", id)
	result, err := r.db.GetDraftByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("draft.get_by_id not found", "id", id)
			return nil, nil
		}
		log.Error("draft.get_by_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("draft.get_by_id succeeded", "id", id)
	return toDraftModelWithRelations(ctx, result), nil
}

func (r *DraftRepository) GetBySubID(ctx context.Context, subID string) (*models.Draft, error) {
	log := observability.NewLogger(ctx, "repository", "draft")
	ctx, span := observability.StartRepoSpan(log.Context(), "draft", "GetBySubID")
	defer span.End()

	log.Info("draft.get_by_sub_id started", "sub_id", subID)
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetDraftBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("draft.get_by_sub_id not found", "sub_id", subID)
			return nil, nil
		}
		log.Error("draft.get_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("draft.get_by_sub_id succeeded", "sub_id", subID)
	return toDraftModelWithRelations(ctx, result), nil
}

func (r *DraftRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	log := observability.NewLogger(ctx, "repository", "draft")
	ctx, span := observability.StartRepoSpan(log.Context(), "draft", "GetIDBySubID")
	defer span.End()

	log.Info("draft.get_id_by_sub_id started", "sub_id", subID)
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetDraftBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("draft.get_id_by_sub_id not found", "sub_id", subID)
			return 0, nil
		}
		log.Error("draft.get_id_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return 0, err
	}
	log.Info("draft.get_id_by_sub_id succeeded", "sub_id", subID, "id", result.ID)
	return result.ID, nil
}

func (r *DraftRepository) UpdateBySubID(ctx context.Context, subID string, params query.UpdateDraftParams) (*models.Draft, error) {
	log := observability.NewLogger(ctx, "repository", "draft")
	ctx, span := observability.StartRepoSpan(log.Context(), "draft", "UpdateBySubID")
	defer span.End()

	log.Info("draft.update_by_sub_id started", "sub_id", subID)
	params.SubID = pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	_, err := r.db.UpdateDraft(ctx, params)
	if err != nil {
		log.Error("draft.update_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		log.Debug("draft.update_by_sub_id not found", "sub_id", subID)
		return nil, nil
	}
	log.Info("draft.update_by_sub_id succeeded", "sub_id", subID)
	return r.GetByID(ctx, id)
}

func (r *DraftRepository) UpdateStatus(ctx context.Context, subID string, status string) error {
	log := observability.NewLogger(ctx, "repository", "draft")
	ctx, span := observability.StartRepoSpan(log.Context(), "draft", "UpdateStatus")
	defer span.End()

	log.Info("draft.update_status started", "sub_id", subID, "status", status)
	params := query.UpdateDraftStatusParams{
		SubID:  pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true},
		Status: status,
	}
	_, err := r.db.UpdateDraftStatus(ctx, params)
	if err == pgx.ErrNoRows {
		log.Debug("draft.update_status not found", "sub_id", subID)
		observability.RecordError(ctx, entities.ErrDraftNotFound)
		return entities.ErrDraftNotFound
	}
	if err != nil {
		log.Error("draft.update_status failed", "error", err)
		observability.RecordError(ctx, err)
	}
	return err
}

func (r *DraftRepository) List(ctx context.Context, params query.ListDraftsParams) ([]*models.Draft, int64, error) {
	log := observability.NewLogger(ctx, "repository", "draft")
	ctx, span := observability.StartRepoSpan(log.Context(), "draft", "List")
	defer span.End()

	log.Info("draft.list started")
	rows, err := r.db.ListDrafts(ctx, params)
	if err != nil {
		log.Error("draft.list failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, 0, err
	}

	drafts := make([]*models.Draft, 0, len(rows))
	var total int64
	for _, row := range rows {
		drafts = append(drafts, toDraftModelWithRelations(ctx, row))
		total = row.Total
	}

	log.Info("draft.list succeeded", "count", len(drafts), "total", total)
	return drafts, total, nil
}

func (r *DraftRepository) SoftDelete(ctx context.Context, subID string) error {
	log := observability.NewLogger(ctx, "repository", "draft")
	ctx, span := observability.StartRepoSpan(log.Context(), "draft", "SoftDelete")
	defer span.End()

	log.Info("draft.soft_delete started", "sub_id", subID)
	_, err := r.db.SoftDeleteDraft(ctx, pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true})
	if err != nil {
		log.Error("draft.soft_delete failed", "error", err)
		observability.RecordError(ctx, err)
	}
	return err
}

func toDraftModelWithRelations(ctx context.Context, q interface{}) *models.Draft {
	_, span := observability.StartRepoSpan(ctx, "drafts", "to_model_with_relations")
	defer span.End()

	var (
		subID                pgtype.UUID
		accountSubID         pgtype.UUID
		categorySubID        pgtype.UUID
		transferAccountSubID pgtype.UUID
		transactionType      string
		title                string
		baseAmount           int64
		enhancedAmount       pgtype.Int8
		currency             string
		currencyRate         pgtype.Numeric
		notes                pgtype.Text
		source               string
		status               string
		confirmedAt          pgtype.Timestamptz
		rejectedAt           pgtype.Timestamptz
		createdAt            pgtype.Timestamptz
		updatedAt            pgtype.Timestamptz
		deletedAt            pgtype.Timestamptz
		transferAccountID    pgtype.Int4
	)

	switch v := q.(type) {
	case query.GetDraftByIDRow:
		subID = v.SubID
		accountSubID = v.AccountSubID
		categorySubID = v.CategorySubID
		transferAccountSubID = v.TransferAccountSubID
		transactionType = v.TransactionType
		title = v.Title
		baseAmount = v.BaseAmount
		enhancedAmount = v.EnhancedAmount
		currency = v.Currency
		currencyRate = v.CurrencyRate
		notes = v.Notes
		source = v.Source
		status = v.Status
		confirmedAt = v.ConfirmedAt
		rejectedAt = v.RejectedAt
		createdAt = v.CreatedAt
		updatedAt = v.UpdatedAt
		deletedAt = v.DeletedAt
		transferAccountID = v.TransferAccountID
	case query.GetDraftBySubIDRow:
		subID = v.SubID
		accountSubID = v.AccountSubID
		categorySubID = v.CategorySubID
		transferAccountSubID = v.TransferAccountSubID
		transactionType = v.TransactionType
		title = v.Title
		baseAmount = v.BaseAmount
		enhancedAmount = v.EnhancedAmount
		currency = v.Currency
		currencyRate = v.CurrencyRate
		notes = v.Notes
		source = v.Source
		status = v.Status
		confirmedAt = v.ConfirmedAt
		rejectedAt = v.RejectedAt
		createdAt = v.CreatedAt
		updatedAt = v.UpdatedAt
		deletedAt = v.DeletedAt
		transferAccountID = v.TransferAccountID
	case query.ListDraftsRow:
		subID = v.SubID
		accountSubID = v.AccountSubID
		categorySubID = v.CategorySubID
		transferAccountSubID = v.TransferAccountSubID
		transactionType = v.TransactionType
		title = v.Title
		baseAmount = v.BaseAmount
		enhancedAmount = v.EnhancedAmount
		currency = v.Currency
		currencyRate = v.CurrencyRate
		notes = v.Notes
		source = v.Source
		status = v.Status
		confirmedAt = v.ConfirmedAt
		rejectedAt = v.RejectedAt
		createdAt = v.CreatedAt
		updatedAt = v.UpdatedAt
		deletedAt = v.DeletedAt
		transferAccountID = v.TransferAccountID
	default:
		return nil
	}

	amount := baseAmount
	if enhancedAmount.Valid {
		amount = enhancedAmount.Int64
	}

	var currencyRateFloat float64
	if currencyRate.Valid {
		if fv, err := currencyRate.Float64Value(); err == nil {
			currencyRateFloat = fv.Float64
		}
	}

	m := &models.Draft{
		SubID:           uuid.UUID(subID.Bytes).String(),
		AccountID:       uuid.UUID(accountSubID.Bytes).String(),
		CategoryID:      uuid.UUID(categorySubID.Bytes).String(),
		TransactionType: transactionType,
		Title:           title,
		Amount:          amount,
		Currency:        currency,
		CurrencyRate:    currencyRateFloat,
		Notes:           notes.String,
		Source:          source,
		Status:          status,
		CreatedAt:       createdAt.Time,
		UpdatedAt:       updatedAt.Time,
	}

	if transferAccountID.Valid {
		if transferAccountSubID.Valid {
			m.TransferAccountID = uuid.UUID(transferAccountSubID.Bytes).String()
		}
	}

	if confirmedAt.Valid {
		m.ConfirmedAt = &confirmedAt.Time
	}

	if rejectedAt.Valid {
		m.RejectedAt = &rejectedAt.Time
	}

	if deletedAt.Valid {
		m.DeletedAt = &deletedAt.Time
	}

	return m
}
