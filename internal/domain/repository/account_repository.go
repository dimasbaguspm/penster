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
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type AccountRepository struct {
	db *query.Queries
}

func NewAccountRepository(db *query.Queries) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, params query.CreateAccountParams) (*models.Account, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "Create")
	defer span.End()

	result, err := r.db.CreateAccount(ctx, params)
	if err != nil {
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id int32) (*models.Account, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "GetByID")
	defer span.End()

	result, err := r.db.GetAccountByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) GetBySubID(ctx context.Context, subID string) (*models.Account, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "GetBySubID")
	defer span.End()

	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetAccountBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "GetIDBySubID")
	defer span.End()

	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetAccountBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		observability.RecordError(ctx, err)
		return 0, err
	}
	return result.ID, nil
}

func (r *AccountRepository) UpdateBySubID(ctx context.Context, subID string, params query.UpdateAccountParams) (*models.Account, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "UpdateBySubID")
	defer span.End()

	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Update(ctx, id, params)
}

func (r *AccountRepository) DeleteBySubID(ctx context.Context, subID string) (*models.Account, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "DeleteBySubID")
	defer span.End()

	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Delete(ctx, id)
}

func (r *AccountRepository) List(ctx context.Context, params query.ListAccountsParams) ([]*models.Account, int64, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "List")
	defer span.End()

	rows, err := r.db.ListAccounts(ctx, params)
	if err != nil {
		observability.RecordError(ctx, err)
		return nil, 0, err
	}

	accounts := make([]*models.Account, 0, len(rows))
	var total int64
	for _, row := range rows {
		accounts = append(accounts, toAccountModel(ctx, query.Account{
			ID:        row.ID,
			SubID:     row.SubID,
			Name:      row.Name,
			Type:      row.Type,
			Balance:   row.Balance,
			DeletedAt: row.DeletedAt,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}))
		total = row.Total
	}

	return accounts, total, nil
}

func (r *AccountRepository) Update(ctx context.Context, id int32, params query.UpdateAccountParams) (*models.Account, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "Update")
	defer span.End()

	result, err := r.db.UpdateAccount(ctx, query.UpdateAccountParams{
		Name:    params.Name,
		Type:    params.Type,
		Balance: params.Balance,
		ID:      id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) Delete(ctx context.Context, id int32) (*models.Account, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "Delete")
	defer span.End()

	result, err := r.db.DeleteAccount(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) UpdateBalanceByID(ctx context.Context, id int32, newBalance int64) (*models.Account, error) {
	ctx, span := observability.StartRepoSpan(ctx, "accounts", "UpdateBalanceByID")
	defer span.End()

	result, err := r.db.UpdateAccountBalance(ctx, query.UpdateAccountBalanceParams{
		Balance: newBalance,
		ID:      id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		observability.RecordError(ctx, err)
		return nil, err
	}
	return toAccountModel(ctx, result), nil
}

func toAccountModel(ctx context.Context, q query.Account) *models.Account {
	_, span := observability.StartRepoSpan(ctx, "accounts", "to_model")
	defer span.End()

	m := &models.Account{
		SubID:     uuid.UUID(q.SubID.Bytes).String(),
		Name:      q.Name,
		Balance:   q.Balance,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	m.Type = models.AccountType(q.Type)

	if q.DeletedAt.Valid {
		m.DeletedAt = &q.DeletedAt.Time
	}
	if q.CreatedAt.Valid {
		m.CreatedAt = q.CreatedAt.Time
	}
	if q.UpdatedAt.Valid {
		m.UpdatedAt = q.UpdatedAt.Time
	}

	return m
}
