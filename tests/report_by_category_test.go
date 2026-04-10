package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestReportByCategory_EmptyPeriod verifies by-category report with no transactions.
func TestReportByCategory_EmptyPeriod(t *testing.T) {
	// Use raw HTTP to see the actual error
	resp, err := http.Get(serverURL + "/reports/by-category?start_date=2099-01-01&end_date=2099-01-31")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
	t.Logf("Status: %d", resp.StatusCode)
	// Read body
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response body: %s", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestReportByCategory_WithTransactions verifies by-category report with transactions.
func TestReportByCategory_WithTransactions(t *testing.T) {
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
		Title:           "ByCategory Test Expense",
		Amount:          400,
		Currency:        "USD",
	}
	doCreateTransaction(expenseReq)

	// Query report - use current year since transactions are created with time.Now()
	_, status, err := doGetReportByCategory(yearStr+"-01-01", yearStr+"-12-31")
	if err != nil {
		t.Fatalf("Failed to get report by category: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
}

// TestReportByCategory_LargeRange verifies by-category report with 1-year range.
func TestReportByCategory_LargeRange(t *testing.T) {
	year := time.Now().Year()
	yearStr := fmt.Sprintf("%d", year)

	account := createTestAccountWithBalance(t, 1000)
	category := createTestDraftCategory(t)

	// Create transaction
	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Large Range Test",
		Amount:          150,
		Currency:        "USD",
	}
	doCreateTransaction(req)

	result, status, err := doGetReportByCategory(yearStr+"-01-01", yearStr+"-12-31")
	if err != nil {
		t.Fatalf("Failed to get report by category: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	_ = result
}

// TestReportByCategory_MultipleCategories verifies by-category with multiple categories.
func TestReportByCategory_MultipleCategories(t *testing.T) {
	year := time.Now().Year()
	yearStr := fmt.Sprintf("%d", year)

	account := createTestAccountWithBalance(t, 1000)
	category1 := createTestDraftCategory(t)
	category2 := createTestDraftCategory(t)

	// Create expense for category1
	expense1 := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category1.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Category1 Expense",
		Amount:          100,
		Currency:        "USD",
	}
	doCreateTransaction(expense1)

	// Create expense for category2
	expense2 := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category2.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Category2 Expense",
		Amount:          200,
		Currency:        "USD",
	}
	doCreateTransaction(expense2)

	// Query report - use current year since transactions are created with time.Now()
	_, status, err := doGetReportByCategory(yearStr+"-01-01", yearStr+"-12-31")
	if err != nil {
		t.Fatalf("Failed to get report by category: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
}
