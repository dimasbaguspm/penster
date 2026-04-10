package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestReportTrends_WithTransactions verifies trends report with transactions.
func TestReportTrends_WithTransactions(t *testing.T) {
	year := time.Now().Year()
	yearStr := fmt.Sprintf("%d", year)

	// Create test data
	account := createTestAccountWithBalance(t, 1000)
	category := createTestDraftCategory(t)

	// Create expense transaction
	expenseReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Trends Test Expense",
		Amount:          200,
		Currency:        "USD",
	}
	doCreateTransaction(expenseReq)

	// Query report - use current year since transactions are created with time.Now()
	result, status, err := doGetReportTrends(yearStr+"-01-01", yearStr+"-12-31")
	if err != nil {
		t.Fatalf("Failed to get report trends: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	_ = result.DataPoints
}

// TestReportTrends_EmptyPeriod verifies trends report with no transactions.
func TestReportTrends_EmptyPeriod(t *testing.T) {
	result, status, err := doGetReportTrends("2099-01-01", "2099-01-31")
	if err != nil {
		t.Fatalf("Failed to get report trends: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	_ = result
}

// TestReportTrends_YearBoundary verifies trends report spanning year boundary.
func TestReportTrends_YearBoundary(t *testing.T) {
	_, status, err := doGetReportTrends("2023-12-31", "2024-01-01")
	if err != nil {
		t.Fatalf("Failed to get report trends: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
}

// TestReportTrends_MultipleDays verifies trends report with multiple days.
func TestReportTrends_MultipleDays(t *testing.T) {
	year := time.Now().Year()
	yearStr := fmt.Sprintf("%d", year)

	_, status, err := doGetReportTrends(yearStr+"-01-01", yearStr+"-12-31")
	if err != nil {
		t.Fatalf("Failed to get report trends: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
}