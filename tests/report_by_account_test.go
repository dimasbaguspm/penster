package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestReportByAccount_WithTransactions verifies by-account report with transactions.
func TestReportByAccount_WithTransactions(t *testing.T) {
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
		Title:           "ByAccount Test Expense",
		Amount:          300,
		Currency:        "USD",
	}
	doCreateTransaction(expenseReq)

	// Query report - use current year since transactions are created with time.Now()
	_, status, err := doGetReportByAccount(yearStr+"-01-01", yearStr+"-12-31")
	if err != nil {
		t.Fatalf("Failed to get report by account: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
}

// TestReportByAccount_EmptyPeriod verifies by-account report with no transactions.
func TestReportByAccount_EmptyPeriod(t *testing.T) {
	result, status, err := doGetReportByAccount("2099-01-01", "2099-01-31")
	if err != nil {
		t.Fatalf("Failed to get report by account: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if len(result.Accounts) != 0 {
		t.Errorf("Expected nil or empty accounts for empty period, got %d", len(result.Accounts))
	}
}

// TestReportByAccount_FutureDates verifies by-account report with future dates.
func TestReportByAccount_FutureDates(t *testing.T) {
	result, status, err := doGetReportByAccount("2099-06-01", "2099-06-30")
	if err != nil {
		t.Fatalf("Failed to get report by account: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if len(result.Accounts) != 0 {
		t.Errorf("Expected empty accounts for future dates, got %d", len(result.Accounts))
	}
}

// TestReportByAccount_MultipleTransactions verifies by-account with multiple transactions.
func TestReportByAccount_MultipleTransactions(t *testing.T) {
	year := time.Now().Year()
	yearStr := fmt.Sprintf("%d", year)

	account1 := createTestAccountWithBalance(t, 1000)
	account2 := createTestAccountWithBalance(t, 500)
	category := createTestDraftCategory(t)

	// Create expense for account1
	expense1 := &models.CreateTransactionRequest{
		AccountID:       account1.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Account1 Expense",
		Amount:          100,
		Currency:        "USD",
	}
	doCreateTransaction(expense1)

	// Create expense for account2
	expense2 := &models.CreateTransactionRequest{
		AccountID:       account2.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Account2 Expense",
		Amount:          200,
		Currency:        "USD",
	}
	doCreateTransaction(expense2)

	// Query report - use current year since transactions are created with time.Now()
	_, status, err := doGetReportByAccount(yearStr+"-01-01", yearStr+"-12-31")
	if err != nil {
		t.Fatalf("Failed to get report by account: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
}
