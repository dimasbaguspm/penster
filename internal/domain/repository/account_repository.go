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

type AccountRepository struct {
	db *query.Queries
}

func NewAccountRepository(db *query.Queries) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, req *models.CreateAccountRequest) (*models.Account, error) {
	result, err := r.db.CreateAccount(ctx, query.CreateAccountParams{
		Name:    req.Name,
		Type:    string(req.Type),
		Balance: req.Balance,
	})
	if err != nil {
		return nil, err
	}
	return toAccountModel(result), nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id int32) (*models.Account, error) {
	result, err := r.db.GetAccountByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toAccountModel(result), nil
}

func (r *AccountRepository) GetBySubID(ctx context.Context, subID string) (*models.Account, error) {
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetAccountBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toAccountModel(result), nil
}

func (r *AccountRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetAccountBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return result.ID, nil
}

func (r *AccountRepository) UpdateBySubID(ctx context.Context, subID string, req *models.UpdateAccountRequest) (*models.Account, error) {
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Update(ctx, id, req)
}

func (r *AccountRepository) DeleteBySubID(ctx context.Context, subID string) (*models.Account, error) {
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil
	}
	return r.Delete(ctx, id)
}

func (r *AccountRepository) List(ctx context.Context, params *models.AccountSearchParams) ([]*models.Account, int64, error) {
	queryParams := params.ToQueryParams()
	rows, err := r.db.ListAccounts(ctx, queryParams)
	if err != nil {
		return nil, 0, err
	}

	accounts := make([]*models.Account, 0, len(rows))
	var total int64
	for _, row := range rows {
		accounts = append(accounts, toAccountModel(query.Account{
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

func (r *AccountRepository) Update(ctx context.Context, id int32, req *models.UpdateAccountRequest) (*models.Account, error) {
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

	result, err := r.db.UpdateAccount(ctx, query.UpdateAccountParams{
		Name:    name,
		Type:    accType,
		Balance: balance,
		ID:      id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toAccountModel(result), nil
}

func (r *AccountRepository) Delete(ctx context.Context, id int32) (*models.Account, error) {
	result, err := r.db.DeleteAccount(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toAccountModel(result), nil
}

func (r *AccountRepository) UpdateBalanceByID(ctx context.Context, id int32, newBalance int64) (*models.Account, error) {
	result, err := r.db.UpdateAccountBalance(ctx, query.UpdateAccountBalanceParams{
		Balance: newBalance,
		ID:      id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toAccountModel(result), nil
}

func toAccountModel(q query.Account) *models.Account {
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
