package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/response"
)

// doGetReportSummary GETs /reports/summary with date range and returns ReportSummary + status.
func doGetReportSummary(startDate, endDate string) (*models.ReportSummary, int, error) {
	path := fmt.Sprintf("/reports/summary?start_date=%s&end_date=%s", startDate, endDate)

	req, err := http.NewRequest("GET", serverURL+path, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	var respWrapper response.Response
	if err := json.Unmarshal(body, &respWrapper); err != nil {
		return nil, resp.StatusCode, err
	}

	if !respWrapper.Success {
		return nil, resp.StatusCode, fmt.Errorf("API error: %s", respWrapper.Error)
	}

	dataBytes, err := json.Marshal(respWrapper.Data)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	result := &models.ReportSummary{}
	if err := json.Unmarshal(dataBytes, result); err != nil {
		return nil, resp.StatusCode, err
	}
	return result, resp.StatusCode, nil
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