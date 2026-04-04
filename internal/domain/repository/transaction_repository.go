package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
)

type TransactionRepository struct {
	db *query.Queries
}

func NewTransactionRepository(db *query.Queries) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, params query.CreateTransactionParams) (*models.Transaction, error) {
	id, err := r.db.CreateTransaction(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int32) (*models.Transaction, error) {
	result, err := r.db.GetTransactionByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toTransactionModelWithRelations(result), nil
}

func (r *TransactionRepository) GetBySubID(ctx context.Context, subID string) (*models.Transaction, error) {
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetTransactionBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toTransactionModelWithRelations(result), nil
}

func (r *TransactionRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetTransactionBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return result.ID, nil
}

func (r *TransactionRepository) UpdateBySubID(ctx context.Context, subID string, params query.UpdateTransactionParams) (*models.Transaction, error) {
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	params.ID = id
	_, err = r.db.UpdateTransaction(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *TransactionRepository) DeleteBySubID(ctx context.Context, subID string) (*models.Transaction, error) {
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Delete(ctx, id)
}

func (r *TransactionRepository) List(ctx context.Context, params query.ListTransactionsParams) ([]*models.Transaction, int64, error) {
	rows, err := r.db.ListTransactions(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	transactions := make([]*models.Transaction, 0, len(rows))
	var total int64
	for _, row := range rows {
		transactions = append(transactions, toTransactionModelWithRelations(row))
		total = row.Total
	}

	return transactions, total, nil
}

func (r *TransactionRepository) Update(ctx context.Context, id int32, params query.UpdateTransactionParams) (*models.Transaction, error) {
	params.ID = id
	_, err := r.db.UpdateTransaction(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *TransactionRepository) Delete(ctx context.Context, id int32) (*models.Transaction, error) {
	tx, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	_, err = r.db.DeleteTransaction(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return tx, nil
}

func toTransactionModelWithRelations(q interface{}) *models.Transaction {
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
		createdAt            pgtype.Timestamptz
		updatedAt            pgtype.Timestamptz
		deletedAt            pgtype.Timestamptz
		transferAccountID    pgtype.Int4
	)

	switch v := q.(type) {
	case query.GetTransactionByIDRow:
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
		createdAt = v.CreatedAt
		updatedAt = v.UpdatedAt
		deletedAt = v.DeletedAt
		transferAccountID = v.TransferAccountID
	case query.GetTransactionBySubIDRow:
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
		createdAt = v.CreatedAt
		updatedAt = v.UpdatedAt
		deletedAt = v.DeletedAt
		transferAccountID = v.TransferAccountID
	case query.ListTransactionsRow:
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

	m := &models.Transaction{
		SubID:           uuid.UUID(subID.Bytes).String(),
		AccountID:       uuid.UUID(accountSubID.Bytes).String(),
		CategoryID:      uuid.UUID(categorySubID.Bytes).String(),
		TransactionType: models.TransactionType(transactionType),
		Title:           title,
		Amount:          amount,
		Currency:        currency,
		CurrencyRate:    currencyRateFloat,
		Notes:           notes.String,
		CreatedAt:       createdAt.Time,
		UpdatedAt:       updatedAt.Time,
	}

	if transferAccountID.Valid {
		if transferAccountSubID.Valid {
			m.TransferAccountID = uuid.UUID(transferAccountSubID.Bytes).String()
		}
	}

	if deletedAt.Valid {
		m.DeletedAt = &deletedAt.Time
	}

	return m
}

func ptrString(s string) *string {
	return &s
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
