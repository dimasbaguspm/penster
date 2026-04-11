package valueobjects

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

// ToCreateDraftParams converts CreateDraftRequest with translated IDs to query params
func ToCreateDraftParams(
	ctx context.Context,
	accountID int32,
	transferAccountID pgtype.Int4,
	categoryID int32,
	currencyRate float64,
	req *models.CreateDraftRequest,
) (query.CreateDraftParams, error) {
	ctx, span := observability.StartValueObjectSpan(ctx, "draft", "to_create_params")
	defer span.End()

	enhancedAmountInt := req.Amount * int64(currencyRate)
	var currencyRateNumeric pgtype.Numeric
	if err := currencyRateNumeric.Scan(fmt.Sprintf("%.6f", currencyRate)); err != nil {
		return query.CreateDraftParams{}, err
	}

	return query.CreateDraftParams{
		AccountID:         accountID,
		TransferAccountID: transferAccountID,
		CategoryID:        pgtype.Int4{Int32: categoryID, Valid: true},
		TransactionType:   req.TransactionType,
		Title:             req.Title,
		BaseAmount:        req.Amount,
		EnhancedAmount:    pgtype.Int8{Int64: enhancedAmountInt, Valid: true},
		Currency:          req.Currency,
		CurrencyRate:      currencyRateNumeric,
		Notes:             pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
		Source:            req.Source,
		Status:            string(models.DraftStatusPending),
	}, nil
}

// ToUpdateDraftParams converts UpdateDraftRequest with translated IDs to query params
func ToUpdateDraftParams(
	ctx context.Context,
	accountID int32,
	transferAccountID pgtype.Int4,
	categoryID int32,
	currencyRate float64,
	req *models.UpdateDraftRequest,
) query.UpdateDraftParams {
	ctx, span := observability.StartValueObjectSpan(ctx, "draft", "to_update_params")
	defer span.End()

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
	params := query.UpdateDraftParams{
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
		params.TransactionType = *req.TransactionType
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

// ToListDraftsParams converts DraftSearchParams to query params
func ToListDraftsParams(ctx context.Context, params *models.DraftSearchParams) query.ListDraftsParams {
	ctx, span := observability.StartValueObjectSpan(ctx, "draft", "to_list_params")
	defer span.End()

	queryParams := query.ListDraftsParams{
		Column1:  "",
		Column2:  "",
		PageSize: params.PageSize,
	}

	if params.Source != nil {
		queryParams.Column1 = *params.Source
	}
	if params.Status != nil {
		queryParams.Column2 = *params.Status
	}

	return queryParams
}
