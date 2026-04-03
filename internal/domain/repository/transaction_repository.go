package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/syncerr"
)

type TransactionRepository struct {
	db           *query.Queries
	accountRepo  *AccountRepository
	categoryRepo *CategoryRepository
}

func NewTransactionRepository(db *query.Queries, accountRepo *AccountRepository, categoryRepo *CategoryRepository) *TransactionRepository {
	return &TransactionRepository{
		db:           db,
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, req *models.CreateTransactionRequest, currencyRate float64) (*models.Transaction, error) {
	var (
		accountID         int32
		transferAccountID pgtype.Int4
		categoryID        pgtype.Int4
		enhancedAmount    pgtype.Int8
	)

	grp := syncerr.Group{}

	grp.Go(func() error {
		accID, err := r.accountRepo.GetIDBySubID(ctx, req.AccountID)
		if err != nil {
			return err
		}
		if accID == 0 {
			return fmt.Errorf("account not found: %s", req.AccountID)
		}
		accountID = accID
		return nil
	})

	if req.TransferAccountID != "" {
		grp.Go(func() error {
			transferID, err := r.accountRepo.GetIDBySubID(ctx, req.TransferAccountID)
			if err != nil {
				return err
			}
			if transferID > 0 {
				transferAccountID = pgtype.Int4{Int32: transferID, Valid: true}
			}
			return nil
		})
	}

	grp.Go(func() error {
		catID, err := r.categoryRepo.GetIDBySubID(ctx, req.CategoryID)
		if err != nil {
			return err
		}
		if catID == 0 {
			return fmt.Errorf("category not found: %s", req.CategoryID)
		}
		categoryID = pgtype.Int4{Int32: catID, Valid: true}
		return nil
	})

	if errs := grp.Wait(); len(errs) > 0 {
		return nil, errs[0]
	}

	if currencyRate == 0 {
		currencyRate = 1
	}

	enhancedAmountInt := req.Amount * int64(currencyRate)
	enhancedAmount = pgtype.Int8{Int64: enhancedAmountInt, Valid: true}

	transactedAt := time.Now()

	id, err := r.db.CreateTransaction(ctx, query.CreateTransactionParams{
		AccountID:         accountID,
		TransferAccountID: transferAccountID,
		CategoryID:        categoryID,
		TransactionType:   string(req.TransactionType),
		Title:             req.Title,
		BaseAmount:        req.Amount,
		EnhancedAmount:    enhancedAmount,
		Currency:          req.Currency,
		CurrencyRate:      currencyRate,
		TransactedAt:      pgtype.Date{Time: transactedAt, Valid: true},
		Notes:             pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})
	if err != nil {
		return nil, err
	}

	// Re-query with relations
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

func (r *TransactionRepository) UpdateBySubID(ctx context.Context, subID string, req *models.UpdateTransactionRequest, currencyRate float64) (*models.Transaction, error) {
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Update(ctx, id, req, currencyRate)
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

func (r *TransactionRepository) List(ctx context.Context, params *models.TransactionSearchParams) ([]*models.Transaction, int64, error) {
	var accountIDs []int32
	var categoryIDs []int32

	grp := syncerr.Group{}

	if len(params.AccountIDs) > 0 {
		grp.Go(func() error {
			for _, accSubID := range params.AccountIDs {
				accID, err := r.accountRepo.GetIDBySubID(ctx, accSubID)
				if err != nil {
					return err
				}
				if accID > 0 {
					accountIDs = append(accountIDs, accID)
				}
			}
			return nil
		})
	}

	if len(params.CategoryIDs) > 0 {
		grp.Go(func() error {
			for _, catSubID := range params.CategoryIDs {
				catID, err := r.categoryRepo.GetIDBySubID(ctx, catSubID)
				if err != nil {
					return err
				}
				if catID > 0 {
					categoryIDs = append(categoryIDs, catID)
				}
			}
			return nil
		})
	}

	if errs := grp.Wait(); len(errs) > 0 {
		return nil, 0, errs[0]
	}

	var queryParams query.ListTransactionsParams
	if len(accountIDs) > 0 {
		queryParams.AccountID = pgtype.Int4{Int32: accountIDs[0], Valid: true}
	}
	if len(categoryIDs) > 0 {
		queryParams.CategoryID = pgtype.Int4{Int32: categoryIDs[0], Valid: true}
	}
	if params.Q != nil {
		queryParams.Q = *params.Q
	}
	queryParams.SortBy = params.SortBy
	queryParams.SortOrder = params.SortOrder
	queryParams.PageSize = int64(params.PageSize)

	rows, err := r.db.ListTransactions(ctx, queryParams)
	if err != nil {
		return nil, 0, err
	}

	transactions := make([]*models.Transaction, 0, len(rows))
	var total int64
	for _, row := range rows {
		transactions = append(transactions, toTransactionModelWithRelations(query.TransactionWithRelations{
			ID:                   row.ID,
			SubID:                row.SubID,
			AccountID:            row.AccountID,
			TransferAccountID:    row.TransferAccountID,
			CategoryID:           row.CategoryID,
			TransactionType:      row.TransactionType,
			Title:                row.Title,
			BaseAmount:           row.BaseAmount,
			EnhancedAmount:       row.EnhancedAmount,
			Currency:             row.Currency,
			CurrencyRate:         row.CurrencyRate,
			TransactedAt:         row.TransactedAt,
			Notes:                row.Notes,
			DeletedAt:            row.DeletedAt,
			CreatedAt:            row.CreatedAt,
			UpdatedAt:            row.UpdatedAt,
			AccountSubID:         row.AccountSubID,
			TransferAccountSubID: row.TransferAccountSubID,
			CategorySubID:        row.CategorySubID,
		}))
		total = row.Total
	}

	return transactions, total, nil
}

func (r *TransactionRepository) Update(ctx context.Context, id int32, req *models.UpdateTransactionRequest, currencyRate float64) (*models.Transaction, error) {
	var (
		accountID         pgtype.Int4
		transferAccountID pgtype.Int4
		categoryID        pgtype.Int4
		transactionType   pgtype.Text
		title             pgtype.Text
		amount            pgtype.Int8
		currency          pgtype.Text
		notes             pgtype.Text
	)

	grp := syncerr.Group{}

	if req.AccountID != nil {
		grp.Go(func() error {
			accID, err := r.accountRepo.GetIDBySubID(ctx, *req.AccountID)
			if err != nil {
				return err
			}
			accountID = pgtype.Int4{Int32: accID, Valid: accID > 0}
			return nil
		})
	}

	if req.TransferAccountID != nil {
		grp.Go(func() error {
			accID, err := r.accountRepo.GetIDBySubID(ctx, *req.TransferAccountID)
			if err != nil {
				return err
			}
			transferAccountID = pgtype.Int4{Int32: accID, Valid: accID > 0}
			return nil
		})
	}

	if req.CategoryID != nil {
		grp.Go(func() error {
			catID, err := r.categoryRepo.GetIDBySubID(ctx, *req.CategoryID)
			if err != nil {
				return err
			}
			categoryID = pgtype.Int4{Int32: catID, Valid: catID > 0}
			return nil
		})
	}

	if errs := grp.Wait(); len(errs) > 0 {
		return nil, errs[0]
	}

	if req.TransactionType != nil {
		transactionType = pgtype.Text{String: string(*req.TransactionType), Valid: true}
	}

	if req.Title != nil {
		title = pgtype.Text{String: *req.Title, Valid: true}
	}

	if req.Amount != nil {
		amount = pgtype.Int8{Int64: *req.Amount, Valid: true}
	}

	if req.Currency != nil {
		currency = pgtype.Text{String: *req.Currency, Valid: true}
	}

	if req.Notes != nil {
		notes = pgtype.Text{String: *req.Notes, Valid: true}
	}

	var enhancedAmount pgtype.Int8
	if currencyRate > 0 && amount.Valid {
		enhancedAmountInt := amount.Int64 * int64(currencyRate)
		enhancedAmount = pgtype.Int8{Int64: enhancedAmountInt, Valid: true}
	}

	_, err := r.db.UpdateTransaction(ctx, query.UpdateTransactionParams{
		AccountID:         accountID,
		TransferAccountID: transferAccountID,
		CategoryID:        categoryID,
		TransactionType:   transactionType,
		Title:             title,
		BaseAmount:        amount,
		EnhancedAmount:    enhancedAmount,
		Currency:          currency,
		CurrencyRate:      currencyRate,
		Notes:             notes,
		ID:                id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
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

func toTransactionModelWithRelations(q query.TransactionWithRelations) *models.Transaction {
	amount := q.BaseAmount
	if q.EnhancedAmount.Valid {
		amount = q.EnhancedAmount.Int64
	}

	m := &models.Transaction{
		SubID:           uuid.UUID(q.SubID.Bytes).String(),
		AccountID:       uuid.UUID(q.AccountSubID.Bytes).String(),
		CategoryID:      uuid.UUID(q.CategorySubID.Bytes).String(),
		TransactionType: models.TransactionType(q.TransactionType),
		Title:           q.Title,
		Amount:          amount,
		Currency:        q.Currency,
		CurrencyRate:    q.CurrencyRate,
		Notes:           q.Notes.String,
		CreatedAt:       q.CreatedAt.Time,
		UpdatedAt:       q.UpdatedAt.Time,
	}

	if q.TransferAccountID.Valid {
		if q.TransferAccountSubID.Valid {
			m.TransferAccountID = uuid.UUID(q.TransferAccountSubID.Bytes).String()
		}
	}

	if q.DeletedAt.Valid {
		m.DeletedAt = &q.DeletedAt.Time
	}

	return m
}

func ptrString(s string) *string {
	return &s
}
