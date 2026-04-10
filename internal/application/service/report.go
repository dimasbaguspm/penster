package service

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type ReportService struct {
	query query.ReportQueryInterface
}

func NewReportService(reportQuery query.ReportQueryInterface) *ReportService {
	return &ReportService{
		query: reportQuery,
	}
}

func (s *ReportService) GetSummary(ctx context.Context, startDateStr, endDateStr string) (*models.ReportSummary, error) {
	ctx, span := observability.StartServiceSpan(ctx, "ReportService", "GetSummary")
	defer span.End()

	startDate, endDate, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}

	return s.query.GetReportSummary(ctx, startDate, endDate)
}

func (s *ReportService) GetByAccount(ctx context.Context, startDateStr, endDateStr string) (*models.ReportByAccount, error) {
	ctx, span := observability.StartServiceSpan(ctx, "ReportService", "GetByAccount")
	defer span.End()

	startDate, endDate, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}

	return s.query.GetReportByAccount(ctx, startDate, endDate)
}

func (s *ReportService) GetByCategory(ctx context.Context, startDateStr, endDateStr string) (*models.ReportByCategory, error) {
	ctx, span := observability.StartServiceSpan(ctx, "ReportService", "GetByCategory")
	defer span.End()

	startDate, endDate, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}

	return s.query.GetReportByCategory(ctx, startDate, endDate)
}

func (s *ReportService) GetTrends(ctx context.Context, startDateStr, endDateStr string) (*models.ReportTrends, error) {
	ctx, span := observability.StartServiceSpan(ctx, "ReportService", "GetTrends")
	defer span.End()

	startDate, endDate, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}

	return s.query.GetReportTrends(ctx, startDate, endDate)
}

func parseDates(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, entities.ErrInvalidDateRange
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, entities.ErrInvalidDateRange
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	if startDate.After(endDate) {
		return time.Time{}, time.Time{}, entities.ErrInvalidDateRange
	}

	return startDate, endDate, nil
}
