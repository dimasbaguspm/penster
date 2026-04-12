package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type TransactionRepository struct {
	db *query.Queries
}

func NewTransactionRepository(db *query.Queries) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, params query.CreateTransactionParams) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "repository", "transaction")
	ctx, span := observability.StartRepoSpan(log.Context(), "transaction", "Create")
	defer span.End()

	log.Info("transaction.create started", "account_id", params.AccountID, "category_id", params.CategoryID, "amount", params.BaseAmount)
	id, err := r.db.CreateTransaction(ctx, params)
	if err != nil {
		log.Error("transaction.create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("transaction.created", "id", id)
	return r.GetByID(ctx, id)
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int32) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "repository", "transaction")
	ctx, span := observability.StartRepoSpan(log.Context(), "transaction", "GetByID")
	defer span.End()

	log.Info("transaction.get_by_id started", "id", id)
	result, err := r.db.GetTransactionByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("transaction.get_by_id not found", "id", id)
			return nil, nil
		}
		log.Error("transaction.get_by_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("transaction.get_by_id succeeded", "id", id)
	return toTransactionModelWithRelations(ctx, result), nil
}

func (r *TransactionRepository) GetBySubID(ctx context.Context, subID string) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "repository", "transaction")
	ctx, span := observability.StartRepoSpan(log.Context(), "transaction", "GetBySubID")
	defer span.End()

	log.Info("transaction.get_by_sub_id started", "sub_id", subID)
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetTransactionBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("transaction.get_by_sub_id not found", "sub_id", subID)
			return nil, nil
		}
		log.Error("transaction.get_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("transaction.get_by_sub_id succeeded", "sub_id", subID)
	return toTransactionModelWithRelations(ctx, result), nil
}

func (r *TransactionRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	log := observability.NewLogger(ctx, "repository", "transaction")
	ctx, span := observability.StartRepoSpan(log.Context(), "transaction", "GetIDBySubID")
	defer span.End()

	log.Info("transaction.get_id_by_sub_id started", "sub_id", subID)
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetTransactionBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("transaction.get_id_by_sub_id not found", "sub_id", subID)
			return 0, nil
		}
		log.Error("transaction.get_id_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return 0, err
	}
	log.Info("transaction.get_id_by_sub_id succeeded", "sub_id", subID, "id", result.ID)
	return result.ID, nil
}

func (r *TransactionRepository) UpdateBySubID(ctx context.Context, subID string, params query.UpdateTransactionParams) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "repository", "transaction")
	ctx, span := observability.StartRepoSpan(log.Context(), "transaction", "UpdateBySubID")
	defer span.End()

	log.Info("transaction.update_by_sub_id started", "sub_id", subID)
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		log.Debug("transaction.update_by_sub_id not found", "sub_id", subID)
		return nil, nil
	}
	params.ID = id
	_, err = r.db.UpdateTransaction(ctx, params)
	if err != nil {
		log.Error("transaction.update_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("transaction.update_by_sub_id succeeded", "sub_id", subID)
	return r.GetByID(ctx, id)
}

func (r *TransactionRepository) DeleteBySubID(ctx context.Context, subID string) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "repository", "transaction")
	ctx, span := observability.StartRepoSpan(log.Context(), "transaction", "DeleteBySubID")
	defer span.End()

	log.Info("transaction.delete_by_sub_id started", "sub_id", subID)
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		log.Debug("transaction.delete_by_sub_id not found", "sub_id", subID)
		return nil, nil
	}
	return r.Delete(ctx, id)
}

func (r *TransactionRepository) List(ctx context.Context, params query.ListTransactionsParams) ([]*models.Transaction, int64, error) {
	log := observability.NewLogger(ctx, "repository", "transaction")
	ctx, span := observability.StartRepoSpan(log.Context(), "transaction", "List")
	defer span.End()

	log.Info("transaction.list started")
	rows, err := r.db.ListTransactions(ctx, params)
	if err != nil {
		log.Error("transaction.list failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, 0, err
	}

	transactions := make([]*models.Transaction, 0, len(rows))
	var total int64
	for _, row := range rows {
		transactions = append(transactions, toTransactionModelWithRelations(ctx, row))
		total = row.Total
	}

	log.Info("transaction.list succeeded", "count", len(transactions), "total", total)
	return transactions, total, nil
}

func (r *TransactionRepository) Update(ctx context.Context, id int32, params query.UpdateTransactionParams) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "repository", "transaction")
	ctx, span := observability.StartRepoSpan(log.Context(), "transaction", "Update")
	defer span.End()

	log.Info("transaction.update started", "id", id)
	params.ID = id
	_, err := r.db.UpdateTransaction(ctx, params)
	if err != nil {
		log.Error("transaction.update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("transaction.update succeeded", "id", id)
	return r.GetByID(ctx, id)
}

func (r *TransactionRepository) Delete(ctx context.Context, id int32) (*models.Transaction, error) {
	log := observability.NewLogger(ctx, "repository", "transaction")
	ctx, span := observability.StartRepoSpan(log.Context(), "transaction", "Delete")
	defer span.End()

	log.Info("transaction.delete started", "id", id)
	tx, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		log.Debug("transaction.delete not found", "id", id)
		return nil, nil
	}

	_, err = r.db.DeleteTransaction(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("transaction.delete not found", "id", id)
			return nil, nil
		}
		log.Error("transaction.delete failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	log.Info("transaction.deleted", "id", id)
	return tx, nil
}

func toTransactionModelWithRelations(ctx context.Context, q interface{}) *models.Transaction {
	_, span := observability.StartRepoSpan(ctx, "transactions", "to_model_with_relations")
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
