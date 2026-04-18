package main

import (
	"fmt"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// doGetReportSummary GETs /reports/summary with date range and returns ReportSummaryResponse + status.
func doGetReportSummary(startDate, endDate string) (*models.ReportSummaryResponse, int, error) {
	path := fmt.Sprintf("/reports/summary?start_date=%s&end_date=%s", startDate, endDate)
	result, status, err := doJSONRequest[models.ReportSummaryResponse]("GET", path, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, nil
}

// doGetReportByAccount GETs /reports/by-account with date range and returns ReportByAccount + status.
func doGetReportByAccount(startDate, endDate string) (*models.ReportByAccount, int, error) {
	path := fmt.Sprintf("/reports/by-account?start_date=%s&end_date=%s", startDate, endDate)
	result, status, err := doJSONRequest[models.ReportByAccount]("GET", path, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doGetReportByCategory GETs /reports/by-category with date range and returns ReportByCategory + status.
func doGetReportByCategory(startDate, endDate string) (*models.ReportByCategory, int, error) {
	path := fmt.Sprintf("/reports/by-category?start_date=%s&end_date=%s", startDate, endDate)
	result, status, err := doJSONRequest[models.ReportByCategory]("GET", path, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doGetReportTrends GETs /reports/trends with date range and returns ReportTrends + status.
func doGetReportTrends(startDate, endDate string) (*models.ReportTrends, int, error) {
	path := fmt.Sprintf("/reports/trends?start_date=%s&end_date=%s", startDate, endDate)
	result, status, err := doJSONRequest[models.ReportTrends]("GET", path, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}