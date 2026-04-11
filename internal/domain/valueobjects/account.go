package valueobjects

import (
	"context"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
	"github.com/jackc/pgx/v5/pgtype"
)

// ToCreateAccountParams converts CreateAccountRequest to query params
func ToCreateAccountParams(ctx context.Context, req *models.CreateAccountRequest) query.CreateAccountParams {
	_, span := observability.StartValueObjectSpan(ctx, "account", "to_create_params")
	defer span.End()

	return query.CreateAccountParams{
		Name:    req.Name,
		Type:    string(req.Type),
		Balance: req.Balance,
	}
}

// ToUpdateAccountParams converts UpdateAccountRequest to query params
func ToUpdateAccountParams(ctx context.Context, req *models.UpdateAccountRequest) query.UpdateAccountParams {
	_, span := observability.StartValueObjectSpan(ctx, "account", "to_update_params")
	defer span.End()

	name := ""
	if req.Name != nil {
		name = *req.Name
	}
	accType := ""
	if req.Type != nil {
		accType = string(*req.Type)
	}
	balance := int64(0)
	if req.Balance != nil {
		balance = *req.Balance
	}

	return query.UpdateAccountParams{
		Name:    name,
		Type:    accType,
		Balance: balance,
	}
}

// ToListAccountsParams converts AccountSearchParams to query params
func ToListAccountsParams(ctx context.Context, params *models.AccountSearchParams) query.ListAccountsParams {
	_, span := observability.StartValueObjectSpan(ctx, "account", "to_list_params")
	defer span.End()

	var subID pgtype.UUID
	if params.SubID != nil {
		subID = pgtype.UUID{Bytes: conv.ParseUUID(*params.SubID), Valid: true}
	}

	return query.ListAccountsParams{
		SubID:     subID,
		Q:         conv.StringPtrToEmpty(params.Q),
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
		PageSize:  int64(params.PageSize),
	}
}
