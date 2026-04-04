package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/domain/entities"
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
			return fmt.Errorf("%w: %s", entities.ErrAccountNotFound, req.AccountID)
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
			if transferID == 0 {
				return fmt.Errorf("%w: %s", entities.ErrTransferAccountNotFound, req.TransferAccountID)
			}
			transferAccountID = pgtype.Int4{Int32: transferID, Valid: true}
			return nil
		})
	}

	grp.Go(func() error {
		catID, err := r.categoryRepo.GetIDBySubID(ctx, req.CategoryID)
		if err != nil {
			return err
		}
		if catID == 0 {
			return fmt.Errorf("%w: %s", entities.ErrCategoryNotFound, req.CategoryID)
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

	var currencyRateNumeric pgtype.Numeric
	if err := currencyRateNumeric.Scan(fmt.Sprintf("%.6f", currencyRate)); err != nil {
		return nil, err
	}

	id, err := r.db.CreateTransaction(ctx, query.CreateTransactionParams{
		AccountID:         accountID,
		TransferAccountID: transferAccountID,
		CategoryID:        categoryID,
		TransactionType:   string(req.TransactionType),
		Title:             req.Title,
		BaseAmount:        req.Amount,
		EnhancedAmount:    enhancedAmount,
		Currency:          req.Currency,
		CurrencyRate:      currencyRateNumeric,
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
	// Column1 (sub_id filter) must be Invalid to make SQL "IS NULL" check pass when no filter
	queryParams.Column1 = pgtype.UUID{Valid: false}
	var zero int32
	queryParams.Column2 = zero
	queryParams.Column3 = zero
	// Column4 (transaction_type) must be empty string to make SQL "$4 = ''" check pass when no filter
	queryParams.Column4 = ""
	if len(accountIDs) > 0 {
		queryParams.Column2 = accountIDs[0]
	}
	if len(categoryIDs) > 0 {
		queryParams.Column3 = categoryIDs[0]
	}
	if params.Q != nil {
		queryParams.Column5 = *params.Q
	}
	queryParams.Column6 = params.SortBy
	queryParams.Column7 = params.SortOrder
	queryParams.Column12 = params.PageSize

	rows, err := r.db.ListTransactions(ctx, queryParams)
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

func (r *TransactionRepository) Update(ctx context.Context, id int32, req *models.UpdateTransactionRequest, currencyRate float64) (*models.Transaction, error) {
	var (
		accountID          int32
		transferAccountID  pgtype.Int4
		categoryID         int32
		transactionType    string
		title              string
		amount             int64
		currency           string
		notes              pgtype.Text
		hasAccountID       bool
		hasCategoryID      bool
		hasTransactionType bool
		hasTitle           bool
		hasAmount          bool
		hasCurrency        bool
	)

	grp := syncerr.Group{}

	if req.AccountID != nil {
		grp.Go(func() error {
			accID, err := r.accountRepo.GetIDBySubID(ctx, *req.AccountID)
			if err != nil {
				return err
			}
			if accID == 0 {
				return fmt.Errorf("%w: %s", entities.ErrAccountNotFound, *req.AccountID)
			}
			accountID = accID
			hasAccountID = true
			return nil
		})
	}

	if req.TransferAccountID != nil {
		grp.Go(func() error {
			accID, err := r.accountRepo.GetIDBySubID(ctx, *req.TransferAccountID)
			if err != nil {
				return err
			}
			if accID == 0 {
				return fmt.Errorf("%w: %s", entities.ErrTransferAccountNotFound, *req.TransferAccountID)
			}
			transferAccountID = pgtype.Int4{Int32: accID, Valid: true}
			return nil
		})
	}

	if req.CategoryID != nil {
		grp.Go(func() error {
			catID, err := r.categoryRepo.GetIDBySubID(ctx, *req.CategoryID)
			if err != nil {
				return err
			}
			if catID == 0 {
				return fmt.Errorf("%w: %s", entities.ErrCategoryNotFound, *req.CategoryID)
			}
			categoryID = catID
			hasCategoryID = true
			return nil
		})
	}

	if errs := grp.Wait(); len(errs) > 0 {
		return nil, errs[0]
	}

	if req.TransactionType != nil {
		transactionType = string(*req.TransactionType)
		hasTransactionType = true
	}

	if req.Title != nil {
		title = *req.Title
		hasTitle = true
	}

	if req.Amount != nil {
		amount = *req.Amount
		hasAmount = true
	}

	if req.Currency != nil {
		currency = *req.Currency
		hasCurrency = true
	}

	if req.Notes != nil {
		notes = pgtype.Text{String: *req.Notes, Valid: true}
	}

	var enhancedAmount pgtype.Int8
	if currencyRate > 0 && hasAmount {
		enhancedAmountInt := amount * int64(currencyRate)
		enhancedAmount = pgtype.Int8{Int64: enhancedAmountInt, Valid: true}
	}

	var accountIDParam int32
	if hasAccountID {
		accountIDParam = accountID
	}
	var categoryIDParam pgtype.Int4
	if hasCategoryID {
		categoryIDParam = pgtype.Int4{Int32: categoryID, Valid: true}
	}
	var transactionTypeParam string
	if hasTransactionType {
		transactionTypeParam = transactionType
	}
	var titleParam string
	if hasTitle {
		titleParam = title
	}
	var amountParam int64
	if hasAmount {
		amountParam = amount
	}
	var currencyParam string
	if hasCurrency {
		currencyParam = currency
	}

	var currencyRateNumeric pgtype.Numeric
	if currencyRate > 0 {
		if err := currencyRateNumeric.Scan(fmt.Sprintf("%.6f", currencyRate)); err != nil {
			return nil, err
		}
	}

	_, err := r.db.UpdateTransaction(ctx, query.UpdateTransactionParams{
		AccountID:         accountIDParam,
		TransferAccountID: transferAccountID,
		CategoryID:        categoryIDParam,
		TransactionType:   transactionTypeParam,
		Title:             titleParam,
		BaseAmount:        amountParam,
		EnhancedAmount:    enhancedAmount,
		Currency:          currencyParam,
		CurrencyRate:      currencyRateNumeric,
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
