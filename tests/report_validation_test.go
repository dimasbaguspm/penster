package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestReportValidation_MissingParams_TableDriven verifies validation errors using table-driven subtests.
func TestReportValidation_MissingParams_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "missing_start_date",
			path:       "/reports/summary?end_date=2024-01-31",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing_end_date",
			path:       "/reports/by-account?start_date=2024-01-01",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing_both_dates",
			path:       "/reports/by-category",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing_both_dates_trends",
			path:       "/reports/trends",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, status, _ := doJSONRequest[models.ErrorResponse]("GET", tt.path, nil)
			if status != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, status)
			}
		})
	}
}

// TestReportValidation_InvalidDateFormat_TableDriven verifies invalid date format errors using table-driven subtests.
func TestReportValidation_InvalidDateFormat_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "invalid_start_date_format_mdash",
			path:       "/reports/summary?start_date=01-01-2024&end_date=2024-01-31",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid_end_date_format_slash",
			path:       "/reports/by-account?start_date=2024-01-01&end_date=2024/01/31",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid_start_date_month_first",
			path:       "/reports/by-category?start_date=13-01-2024&end_date=2024-01-31",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "non_date_string_start",
			path:       "/reports/trends?start_date=abc&end_date=2024-01-31",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, status, _ := doJSONRequest[models.ErrorResponse]("GET", tt.path, nil)
			if status != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, status)
			}
		})
	}
}

// TestReportValidation_DateRangeErrors_TableDriven verifies date range logic errors using table-driven subtests.
func TestReportValidation_DateRangeErrors_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "start_after_end",
			path:       "/reports/summary?start_date=2024-02-01&end_date=2024-01-01",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty_string_start",
			path:       "/reports/by-account?start_date=&end_date=2024-01-31",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty_string_end",
			path:       "/reports/trends?start_date=2024-01-01&end_date=",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, status, _ := doJSONRequest[models.ErrorResponse]("GET", tt.path, nil)
			if status != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, status)
			}
		})
	}
}