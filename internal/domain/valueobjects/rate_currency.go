package valueobjects

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

// ToUpsertRateCurrencyParams converts UpsertRateCurrencyRequest to query params
func ToUpsertRateCurrencyParams(ctx context.Context, req *models.UpsertRateCurrencyRequest) query.UpsertRateCurrencyParams {
	_, span := observability.StartValueObjectSpan(ctx, "rate_currency", "to_upsert_params")
	defer span.End()

	var rateNumeric pgtype.Numeric
	rateNumeric.Scan(fmt.Sprintf("%.6f", req.Rate))

	return query.UpsertRateCurrencyParams{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rateNumeric,
		RateDate:     pgtype.Date{Time: req.RateDate, Valid: true},
	}
}

// ToListRateCurrenciesParams converts RateCurrencySearchParams to query params
func ToListRateCurrenciesParams(ctx context.Context, params *models.RateCurrencySearchParams) query.ListRateCurrenciesParams {
	_, span := observability.StartValueObjectSpan(ctx, "rate_currency", "to_list_params")
	defer span.End()

	var fromCurrency, toCurrency string
	if params.FromCurrency != nil {
		fromCurrency = *params.FromCurrency
	}
	if params.ToCurrency != nil {
		toCurrency = *params.ToCurrency
	}

	return query.ListRateCurrenciesParams{
		Column1: fromCurrency,
		Column2: toCurrency,
		Column3: params.PageSize,
		Offset:  int32(params.Offset()),
	}
}
