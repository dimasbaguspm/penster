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
	log := observability.NewLogger(ctx, "service", "report")
	ctx, span := observability.StartServiceSpan(log.Context(), "report", "GetSummary")
	defer span.End()

	log.Info("get_summary started", "start_date", startDateStr, "end_date", endDateStr)

	startDate, endDate, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		log.Error("get_summary failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	result, err := s.query.GetReportSummary(ctx, startDate, endDate)
	if err != nil {
		log.Error("get_summary failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("get_summary succeeded", "start_date", startDateStr, "end_date", endDateStr)
	return result, nil
}

func (s *ReportService) GetByAccount(ctx context.Context, startDateStr, endDateStr string) (*models.ReportByAccount, error) {
	log := observability.NewLogger(ctx, "service", "report")
	ctx, span := observability.StartServiceSpan(log.Context(), "report", "GetByAccount")
	defer span.End()

	log.Info("get_by_account started", "start_date", startDateStr, "end_date", endDateStr)

	startDate, endDate, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		log.Error("get_by_account failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	result, err := s.query.GetReportByAccount(ctx, startDate, endDate)
	if err != nil {
		log.Error("get_by_account failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("get_by_account succeeded", "start_date", startDateStr, "end_date", endDateStr)
	return result, nil
}

func (s *ReportService) GetByCategory(ctx context.Context, startDateStr, endDateStr string) (*models.ReportByCategory, error) {
	log := observability.NewLogger(ctx, "service", "report")
	ctx, span := observability.StartServiceSpan(log.Context(), "report", "GetByCategory")
	defer span.End()

	log.Info("get_by_category started", "start_date", startDateStr, "end_date", endDateStr)

	startDate, endDate, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		log.Error("get_by_category failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	result, err := s.query.GetReportByCategory(ctx, startDate, endDate)
	if err != nil {
		log.Error("get_by_category failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("get_by_category succeeded", "start_date", startDateStr, "end_date", endDateStr)
	return result, nil
}

func (s *ReportService) GetTrends(ctx context.Context, startDateStr, endDateStr string) (*models.ReportTrends, error) {
	log := observability.NewLogger(ctx, "service", "report")
	ctx, span := observability.StartServiceSpan(log.Context(), "report", "GetTrends")
	defer span.End()

	log.Info("get_trends started", "start_date", startDateStr, "end_date", endDateStr)

	startDate, endDate, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		log.Error("get_trends failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}

	result, err := s.query.GetReportTrends(ctx, startDate, endDate)
	if err != nil {
		log.Error("get_trends failed", "error", err)
		observability.RecordError(ctx, err)
		return nil, err
	}
	log.Info("get_trends succeeded", "start_date", startDateStr, "end_date", endDateStr)
	return result, nil
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
