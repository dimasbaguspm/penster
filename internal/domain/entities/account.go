package entities

import (
	"errors"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/conv"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/jackc/pgx/v5/pgtype"
)

// Account-specific errors
var (
	ErrAccountNotFound = errors.New("account not found")
)

// ToListAccountsParams converts AccountSearchParams to query params
func ToListAccountsParams(params *models.AccountSearchParams) query.ListAccountsParams {
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
