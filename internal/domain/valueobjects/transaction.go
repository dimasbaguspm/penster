package valueobjects

import (
	"fmt"
	"time"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/jackc/pgx/v5/pgtype"
)

// ToCreateTransactionParams converts CreateTransactionRequest with translated IDs to query params
func ToCreateTransactionParams(
	accountID int32,
	transferAccountID pgtype.Int4,
	categoryID int32,
	currencyRate float64,
	req *models.CreateTransactionRequest,
) (query.CreateTransactionParams, error) {
	enhancedAmountInt := req.Amount * int64(currencyRate)
	var currencyRateNumeric pgtype.Numeric
	if err := currencyRateNumeric.Scan(fmt.Sprintf("%.6f", currencyRate)); err != nil {
		return query.CreateTransactionParams{}, err
	}

	return query.CreateTransactionParams{
		AccountID:         accountID,
		TransferAccountID: transferAccountID,
		CategoryID:        pgtype.Int4{Int32: categoryID, Valid: true},
		TransactionType:   string(req.TransactionType),
		Title:             req.Title,
		BaseAmount:        req.Amount,
		EnhancedAmount:    pgtype.Int8{Int64: enhancedAmountInt, Valid: true},
		Currency:          req.Currency,
		CurrencyRate:      currencyRateNumeric,
		TransactedAt:      pgtype.Date{Time: time.Now(), Valid: true},
		Notes:             pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	}, nil
}

// ToUpdateTransactionParams converts UpdateTransactionRequest with translated IDs to query params
func ToUpdateTransactionParams(
	accountID int32,
	transferAccountID pgtype.Int4,
	categoryID int32,
	currencyRate float64,
	req *models.UpdateTransactionRequest,
) query.UpdateTransactionParams {
	var enhancedAmount pgtype.Int8
	if req.Amount != nil && currencyRate > 0 {
		enhancedAmountInt := *req.Amount * int64(currencyRate)
		enhancedAmount = pgtype.Int8{Int64: enhancedAmountInt, Valid: true}
	}

	var currencyRateNumeric pgtype.Numeric
	if currencyRate > 0 {
		_ = currencyRateNumeric.Scan(fmt.Sprintf("%.6f", currencyRate))
	}

	var zero int32
	params := query.UpdateTransactionParams{
		AccountID:         zero,
		TransferAccountID: pgtype.Int4{},
		CategoryID:        zero,
	}

	if req.AccountID != nil {
		params.AccountID = accountID
	}
	if req.TransferAccountID != nil {
		params.TransferAccountID = transferAccountID
	}
	if req.CategoryID != nil {
		params.CategoryID = categoryID
	}
	if req.TransactionType != nil {
		params.TransactionType = string(*req.TransactionType)
	}
	if req.Title != nil {
		params.Title = *req.Title
	}
	if req.Amount != nil {
		params.BaseAmount = *req.Amount
	}
	if req.Currency != nil {
		params.Currency = *req.Currency
	}
	if req.Notes != nil {
		params.Notes = pgtype.Text{String: *req.Notes, Valid: true}
	}
	if enhancedAmount.Valid {
		params.EnhancedAmount = enhancedAmount
	}
	if currencyRateNumeric.Valid {
		params.CurrencyRate = currencyRateNumeric
	}

	return params
}

// ToListTransactionsParams converts TransactionSearchParams to query params
func ToListTransactionsParams(params *models.TransactionSearchParams) query.ListTransactionsParams {
	queryParams := query.ListTransactionsParams{
		Column1: pgtype.UUID{Valid: false},
		Column2: int32(0),
		Column3: int32(0),
		Column4: "",
	}

	if params.Q != nil {
		queryParams.Column5 = *params.Q
	}
	queryParams.Column6 = params.SortBy
	queryParams.Column7 = params.SortOrder
	queryParams.Column12 = params.PageSize

	return queryParams
}
