package query

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type ReportQueryInterface interface {
	GetReportSummary(ctx context.Context, startDate, endDate time.Time) (*models.ReportSummary, error)
	GetReportByAccount(ctx context.Context, startDate, endDate time.Time) (*models.ReportByAccount, error)
	GetReportByCategory(ctx context.Context, startDate, endDate time.Time) (*models.ReportByCategory, error)
	GetReportTrends(ctx context.Context, startDate, endDate time.Time) (*models.ReportTrends, error)
}

var _ ReportQueryInterface = (*ReportQuery)(nil)

type ReportQuery struct {
	repo *repository.ReportRepository
}

func NewReportQuery(repo *repository.ReportRepository) *ReportQuery {
	return &ReportQuery{repo: repo}
}

func (q *ReportQuery) GetReportSummary(ctx context.Context, startDate, endDate time.Time) (*models.ReportSummary, error) {
	ctx, span := observability.StartQuerySpan(ctx, "report", "get_summary")
	defer span.End()
	return q.repo.GetReportSummary(ctx, startDate, endDate)
}

func (q *ReportQuery) GetReportByAccount(ctx context.Context, startDate, endDate time.Time) (*models.ReportByAccount, error) {
	ctx, span := observability.StartQuerySpan(ctx, "report", "get_by_account")
	defer span.End()
	return q.repo.GetReportByAccount(ctx, startDate, endDate)
}

func (q *ReportQuery) GetReportByCategory(ctx context.Context, startDate, endDate time.Time) (*models.ReportByCategory, error) {
	ctx, span := observability.StartQuerySpan(ctx, "report", "get_by_category")
	defer span.End()
	return q.repo.GetReportByCategory(ctx, startDate, endDate)
}

func (q *ReportQuery) GetReportTrends(ctx context.Context, startDate, endDate time.Time) (*models.ReportTrends, error) {
	ctx, span := observability.StartQuerySpan(ctx, "report", "get_trends")
	defer span.End()
	return q.repo.GetReportTrends(ctx, startDate, endDate)
}
