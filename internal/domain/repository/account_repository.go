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
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "Create")
	defer span.End()

	log.Info("account.create started", "name", params.Name)
	result, err := r.db.CreateAccount(ctx, params)
	if err != nil {
		log.Error("account.create failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("account.created", "id", result.ID)
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id int32) (*models.Account, error) {
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "GetByID")
	defer span.End()

	log.Info("account.get_by_id started", "id", id)
	result, err := r.db.GetAccountByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("account.get_by_id not found", "id", id)
			return nil, nil
		}
		log.Error("account.get_by_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("account.get_by_id succeeded", "id", id)
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) GetBySubID(ctx context.Context, subID string) (*models.Account, error) {
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "GetBySubID")
	defer span.End()

	log.Info("account.get_by_sub_id started", "sub_id", subID)
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetAccountBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("account.get_by_sub_id not found", "sub_id", subID)
			return nil, nil
		}
		log.Error("account.get_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("account.get_by_sub_id succeeded", "sub_id", subID)
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) GetIDBySubID(ctx context.Context, subID string) (int32, error) {
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "GetIDBySubID")
	defer span.End()

	log.Info("account.get_id_by_sub_id started", "sub_id", subID)
	uid := pgtype.UUID{Bytes: conv.ParseUUID(subID), Valid: true}
	result, err := r.db.GetAccountBySubID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("account.get_id_by_sub_id not found", "sub_id", subID)
			return 0, nil
		}
		log.Error("account.get_id_by_sub_id failed", "error", err)
		observability.RecordError(ctx, err)
		return 0, err
	}
	log.Info("account.get_id_by_sub_id succeeded", "sub_id", subID, "id", result.ID)
	return result.ID, nil
}

func (r *AccountRepository) UpdateBySubID(ctx context.Context, subID string, params query.UpdateAccountParams) (*models.Account, error) {
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "UpdateBySubID")
	defer span.End()

	log.Info("account.update_by_sub_id started", "sub_id", subID)
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		log.Debug("account.update_by_sub_id not found", "sub_id", subID)
		return nil, nil
	}
	return r.Update(ctx, id, params)
}

func (r *AccountRepository) DeleteBySubID(ctx context.Context, subID string) (*models.Account, error) {
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "DeleteBySubID")
	defer span.End()

	log.Info("account.delete_by_sub_id started", "sub_id", subID)
	id, err := r.GetIDBySubID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		log.Debug("account.delete_by_sub_id not found", "sub_id", subID)
		return nil, nil
	}
	return r.Delete(ctx, id)
}

func (r *AccountRepository) List(ctx context.Context, params query.ListAccountsParams) ([]*models.Account, int64, error) {
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "List")
	defer span.End()

	log.Info("account.list started")
	rows, err := r.db.ListAccounts(ctx, params)
	if err != nil {
		log.Error("account.list failed", "error", err)
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

	log.Info("account.list succeeded", "count", len(accounts), "total", total)
	return accounts, total, nil
}

func (r *AccountRepository) Update(ctx context.Context, id int32, params query.UpdateAccountParams) (*models.Account, error) {
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "Update")
	defer span.End()

	log.Info("account.update started", "id", id)
	result, err := r.db.UpdateAccount(ctx, query.UpdateAccountParams{
		Name:    params.Name,
		Type:    params.Type,
		Balance: params.Balance,
		ID:      id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("account.update not found", "id", id)
			return nil, nil
		}
		log.Error("account.update failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("account.updated", "id", id)
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) Delete(ctx context.Context, id int32) (*models.Account, error) {
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "Delete")
	defer span.End()

	log.Info("account.delete started", "id", id)
	result, err := r.db.DeleteAccount(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("account.delete not found", "id", id)
			return nil, nil
		}
		log.Error("account.delete failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("account.deleted", "id", id)
	return toAccountModel(ctx, result), nil
}

func (r *AccountRepository) UpdateBalanceByID(ctx context.Context, id int32, newBalance int64) (*models.Account, error) {
	log := observability.NewLogger(ctx, "repository", "account")
	ctx, span := observability.StartRepoSpan(log.Context(), "account", "UpdateBalanceByID")
	defer span.End()

	log.Info("account.update_balance_by_id started", "id", id, "new_balance", newBalance)
	result, err := r.db.UpdateAccountBalance(ctx, query.UpdateAccountBalanceParams{
		Balance: newBalance,
		ID:      id,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("account.update_balance_by_id not found", "id", id)
			return nil, nil
		}
		log.Error("account.update_balance_by_id failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("account.update_balance_by_id succeeded", "id", id)
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
