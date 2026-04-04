package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type DraftRepository struct {
	db *query.Queries
}

func NewDraftRepository(db *query.Queries) *DraftRepository {
	return &DraftRepository{db: db}
}

func (r *DraftRepository) Create(ctx context.Context, params query.CreateDraftParams) (*models.Draft, error) {
	id, err := r.db.CreateDraft(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *DraftRepository) GetByID(ctx context.Context, id int32) (*models.Draft, error) {
	result, err := r.db.GetDraftByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toDraftModelWithRelations(result), nil
}

func (r *DraftRepository) GetBySubID(ctx context.Context, subID string) (*models.Draft, error) {
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetDraftBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toDraftModelWithRelations(result), nil
}

func (r *DraftRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetDraftBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return result.ID, nil
}

func (r *DraftRepository) UpdateBySubID(ctx context.Context, subID string, params query.UpdateDraftParams) (*models.Draft, error) {
	params.SubID = pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	_, err := r.db.UpdateDraft(ctx, params)
	if err != nil {
		return nil, err
	}
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.GetByID(ctx, id)
}

func (r *DraftRepository) UpdateStatus(ctx context.Context, subID string, status string) error {
	params := query.UpdateDraftStatusParams{
		SubID:  pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true},
		Status: status,
	}
	_, err := r.db.UpdateDraftStatus(ctx, params)
	return err
}

func (r *DraftRepository) List(ctx context.Context, params query.ListDraftsParams) ([]*models.Draft, int64, error) {
	rows, err := r.db.ListDrafts(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	drafts := make([]*models.Draft, 0, len(rows))
	var total int64
	for _, row := range rows {
		drafts = append(drafts, toDraftModelWithRelations(row))
		total = row.Total
	}

	return drafts, total, nil
}

func (r *DraftRepository) SoftDelete(ctx context.Context, subID string) error {
	_, err := r.db.SoftDeleteDraft(ctx, pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true})
	return err
}

func toDraftModelWithRelations(q interface{}) *models.Draft {
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
		transactedAt         pgtype.Date
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
		transactedAt = v.TransactedAt
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
		transactedAt = v.TransactedAt
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
		transactedAt = v.TransactedAt
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
		TransactedAt:    transactedAt.Time,
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
