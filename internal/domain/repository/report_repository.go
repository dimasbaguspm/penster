package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReportRepository struct {
	db *query.Queries
}

func NewReportRepository(db *query.Queries) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetReportSummary(ctx context.Context, startDate, endDate time.Time) (*models.ReportSummary, error) {
	// Get total balance across all accounts
	totalBalance, err := r.db.GetTotalBalance(ctx)
	if err != nil {
		return nil, err
	}

	// Get totals by transaction type
	totalsByType, err := r.db.GetTotalsByType(ctx, query.GetTotalsByTypeParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// Initialize summary with defaults
	summary := &models.ReportSummary{
		TotalBalance:   totalBalance.(int64),
		TotalExpenses:  0,
		TotalIncome:    0,
		TotalTransfers: 0,
		PeriodStart:    startDate,
		PeriodEnd:      endDate,
	}

	// Process totals by type
	for _, row := range totalsByType {
		switch row.TransactionType {
		case "expense":
			summary.TotalExpenses = row.Total.(int64)
		case "income":
			summary.TotalIncome = row.Total.(int64)
		case "transfer":
			summary.TotalTransfers = row.Total.(int64)
		}
	}

	return summary, nil
}

func (r *ReportRepository) GetReportByCategory(ctx context.Context, startDate, endDate time.Time) (*models.ReportByCategory, error) {
	categoryBreakdownRows, err := r.db.GetCategoryBreakdown(ctx, query.GetCategoryBreakdownParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	report := &models.ReportByCategory{
		Categories:  make([]models.CategoryBreakdown, 0, len(categoryBreakdownRows)),
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}

	for _, row := range categoryBreakdownRows {
		report.Categories = append(report.Categories, models.CategoryBreakdown{
			CategoryName: row.CategoryName.String,
			Type:         row.CategoryType.String,
			Total:        row.Total.(int64),
		})
	}

	return report, nil
}

func (r *ReportRepository) GetReportByAccount(ctx context.Context, startDate, endDate time.Time) (*models.ReportByAccount, error) {
	accountBreakdownRows, err := r.db.GetAccountBreakdown(ctx, query.GetAccountBreakdownParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	report := &models.ReportByAccount{
		Accounts:    make([]models.AccountBreakdown, 0, len(accountBreakdownRows)),
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}

	for _, row := range accountBreakdownRows {
		breakdown := models.AccountBreakdown{
			AccountName: row.AccountName.String,
			Type:        row.TransactionType,
			Total:       row.Total.(int64),
		}
		if row.AccountID.Valid {
			breakdown.AccountID = strconv.FormatInt(int64(row.AccountID.Int32), 10)
		}
		report.Accounts = append(report.Accounts, breakdown)
	}

	return report, nil
}

func (r *ReportRepository) GetReportTrends(ctx context.Context, startDate, endDate time.Time) (*models.ReportTrends, error) {
	trendRows, err := r.db.GetTrends(ctx, query.GetTrendsParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	report := &models.ReportTrends{
		DataPoints:  make([]models.TrendDataPoint, 0, len(trendRows)),
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}

	for _, row := range trendRows {
		report.DataPoints = append(report.DataPoints, models.TrendDataPoint{
			Date:  row.Date.Time.Format("2006-01-02"),
			Type:  row.TransactionType,
			Total: row.Total.(int64),
		})
	}

	return report, nil
}
