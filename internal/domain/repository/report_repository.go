package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReportRepository struct {
	db *query.Queries
}

func NewReportRepository(db *query.Queries) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetReportSummary(ctx context.Context, startDate, endDate time.Time) (*models.ReportSummary, error) {
	log := observability.NewLogger(ctx, "repository", "report")
	ctx, span := observability.StartRepoSpan(log.Context(), "report", "GetReportSummary")
	defer span.End()

	log.Info("report.get_summary started", "start_date", startDate.Format("2006-01-02"), "end_date", endDate.Format("2006-01-02"))

	// Get total balance across all accounts
	totalBalance, err := r.db.GetTotalBalance(ctx)
	if err != nil {
		log.Error("report.get_summary failed: get_total_balance", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	// Get totals by transaction type
	totalsByType, err := r.db.GetTotalsByType(ctx, query.GetTotalsByTypeParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		log.Error("report.get_summary failed: get_totals_by_type", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	// Convert totalBalance from interface{} to int64
	var balanceInt64 int64
	if nb, ok := totalBalance.(pgtype.Numeric); ok {
		if nb.Valid {
			if iv, err := nb.Int64Value(); err == nil && iv.Valid {
				balanceInt64 = iv.Int64
			}
		}
	}

	// Initialize summary with defaults
	summary := &models.ReportSummary{
		TotalBalance:   balanceInt64,
		TotalExpenses:  0,
		TotalIncome:    0,
		TotalTransfers: 0,
		PeriodStart:    startDate,
		PeriodEnd:      endDate,
	}

	// Process totals by type
	for _, row := range totalsByType {
		var totalInt64 int64
		// Handle both pgtype.Numeric (from SUM) and direct int64/int
		switch v := row.Total.(type) {
		case int64:
			totalInt64 = v
		case int:
			totalInt64 = int64(v)
		case pgtype.Numeric:
			if v.Valid && !v.NaN {
				if iv, err := v.Int64Value(); err == nil && iv.Valid {
					totalInt64 = iv.Int64
				}
			}
		case []byte:
			// Handle byte slice (some drivers return numeric as []byte)
			if len(v) > 0 {
				if parsed, err := strconv.ParseInt(string(v), 10, 64); err == nil {
					totalInt64 = parsed
				}
			}
		}
		switch row.TransactionType {
		case "expense":
			summary.TotalExpenses = totalInt64
		case "income":
			summary.TotalIncome = totalInt64
		case "transfer":
			summary.TotalTransfers = totalInt64
		}
	}

	return summary, nil
}

func (r *ReportRepository) GetReportByCategory(ctx context.Context, startDate, endDate time.Time) (*models.ReportByCategory, error) {
	log := observability.NewLogger(ctx, "repository", "report")
	ctx, span := observability.StartRepoSpan(log.Context(), "report", "GetReportByCategory")
	defer span.End()

	log.Info("report.get_by_category started", "start_date", startDate.Format("2006-01-02"), "end_date", endDate.Format("2006-01-02"))

	categoryBreakdownRows, err := r.db.GetCategoryBreakdown(ctx, query.GetCategoryBreakdownParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		log.Error("report.get_by_category failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, fmt.Errorf("GetCategoryBreakdown failed: %w", err)
	}

	report := &models.ReportByCategory{
		Categories:  make([]models.CategoryBreakdown, 0),
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}

	for _, row := range categoryBreakdownRows {
		var totalInt64 int64
		if nb, ok := row.Total.(pgtype.Numeric); ok && nb.Valid {
			if iv, err := nb.Int64Value(); err == nil && iv.Valid {
				totalInt64 = iv.Int64
			}
		}
		catName := ""
		if row.CategoryName.Valid {
			catName = row.CategoryName.String
		}
		catType := ""
		if row.CategoryType.Valid {
			catType = row.CategoryType.String
		}
		report.Categories = append(report.Categories, models.CategoryBreakdown{
			CategoryName: catName,
			Type:         catType,
			Total:        totalInt64,
		})
	}

	return report, nil
}

func (r *ReportRepository) GetReportByAccount(ctx context.Context, startDate, endDate time.Time) (*models.ReportByAccount, error) {
	log := observability.NewLogger(ctx, "repository", "report")
	ctx, span := observability.StartRepoSpan(log.Context(), "report", "GetReportByAccount")
	defer span.End()

	log.Info("report.get_by_account started", "start_date", startDate.Format("2006-01-02"), "end_date", endDate.Format("2006-01-02"))

	accountBreakdownRows, err := r.db.GetAccountBreakdown(ctx, query.GetAccountBreakdownParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		log.Error("report.get_by_account failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	report := &models.ReportByAccount{
		Accounts:    make([]models.AccountBreakdown, 0, len(accountBreakdownRows)),
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}

	for _, row := range accountBreakdownRows {
		var totalInt64 int64
		if nb, ok := row.Total.(pgtype.Numeric); ok {
			if nb.Valid {
				if iv, err := nb.Int64Value(); err == nil && iv.Valid {
					totalInt64 = iv.Int64
				}
			}
		}
		breakdown := models.AccountBreakdown{
			AccountName: row.AccountName.String,
			Type:        row.TransactionType,
			Total:       totalInt64,
		}
		if row.AccountID.Valid {
			breakdown.AccountID = strconv.FormatInt(int64(row.AccountID.Int32), 10)
		}
		report.Accounts = append(report.Accounts, breakdown)
	}

	return report, nil
}

func (r *ReportRepository) GetReportTrends(ctx context.Context, startDate, endDate time.Time) (*models.ReportTrends, error) {
	log := observability.NewLogger(ctx, "repository", "report")
	ctx, span := observability.StartRepoSpan(log.Context(), "report", "GetReportTrends")
	defer span.End()

	log.Info("report.get_trends started", "start_date", startDate.Format("2006-01-02"), "end_date", endDate.Format("2006-01-02"))

	trendRows, err := r.db.GetTrends(ctx, query.GetTrendsParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		log.Error("report.get_trends failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	report := &models.ReportTrends{
		DataPoints:  make([]models.TrendDataPoint, 0, len(trendRows)),
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}

	for _, row := range trendRows {
		var totalInt64 int64
		if nb, ok := row.Total.(pgtype.Numeric); ok {
			if nb.Valid {
				if iv, err := nb.Int64Value(); err == nil && iv.Valid {
					totalInt64 = iv.Int64
				}
			}
		}
		report.DataPoints = append(report.DataPoints, models.TrendDataPoint{
			Date:  row.Date.Time.Format("2006-01-02"),
			Type:  row.TransactionType,
			Total: totalInt64,
		})
	}

	return report, nil
}
