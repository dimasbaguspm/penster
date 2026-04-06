package query

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/pkg/models"
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
	return q.repo.GetReportSummary(ctx, startDate, endDate)
}

func (q *ReportQuery) GetReportByAccount(ctx context.Context, startDate, endDate time.Time) (*models.ReportByAccount, error) {
	return q.repo.GetReportByAccount(ctx, startDate, endDate)
}

func (q *ReportQuery) GetReportByCategory(ctx context.Context, startDate, endDate time.Time) (*models.ReportByCategory, error) {
	return q.repo.GetReportByCategory(ctx, startDate, endDate)
}

func (q *ReportQuery) GetReportTrends(ctx context.Context, startDate, endDate time.Time) (*models.ReportTrends, error) {
	return q.repo.GetReportTrends(ctx, startDate, endDate)
}
