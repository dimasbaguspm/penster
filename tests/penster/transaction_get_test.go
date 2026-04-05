package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestGetTransaction_Success(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Get Test Transaction",
		Amount:          300,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	result, status, err := doGetTransaction(id)
	if err != nil {
		t.Fatalf("Failed to get transaction: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.SubID != id {
		t.Errorf("Expected id '%s', got %s", id, result.Data.SubID)
	}
}

func TestGetTransaction_NotFound(t *testing.T) {
	_, status, _ := doGetTransaction("00000000-0000-0000-0000-000000000000")
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

func TestGetTransaction_InvalidUUID(t *testing.T) {
	_, status, _ := doGetTransaction("not-a-uuid")
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}

func TestListTransactions_Success(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	// Create a few transactions
	for i := 0; i < 3; i++ {
		req := &models.CreateTransactionRequest{
			AccountID:       account.Data.SubID,
			CategoryID:      category.Data.SubID,
			TransactionType: models.TransactionTypeExpense,
			Title:           "List Test Transaction",
			Amount:          int64(100 * (i + 1)),
			Currency:        "USD",
		}
		_, _, _ = doCreateTransaction(req)
	}

	result, status, err := doListTransactions()
	if err != nil {
		t.Fatalf("Failed to list transactions: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data == nil {
		t.Errorf("Expected data to be non-nil")
	}
	if result.Meta == nil {
		t.Errorf("Expected meta to be non-nil")
	}
}

func TestListTransactions_PaginationMeta(t *testing.T) {
	result, status, err := doListTransactions()
	if err != nil {
		t.Fatalf("Failed to list transactions: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Meta == nil {
		t.Fatalf("Expected meta to be non-nil")
	}
	if result.Meta.Page <= 0 {
		t.Errorf("Expected page >= 1, got %d", result.Meta.Page)
	}
	if result.Meta.PerPage <= 0 {
		t.Errorf("Expected per_page >= 1, got %d", result.Meta.PerPage)
	}
	if result.Meta.Total < 0 {
		t.Errorf("Expected total >= 0, got %d", result.Meta.Total)
	}
}

func TestListTransactions_ReturnsCreatedTransactions(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeIncome,
		Title:           "Specific Income Transaction",
		Amount:          999,
		Currency:        "EUR",
	}
	created, createStatus, createErr := doCreateTransaction(req)
	if createErr != nil {
		t.Fatalf("Failed to create transaction: %v", createErr)
	}
	if createStatus != http.StatusCreated {
		t.Fatalf("Failed to create transaction: status %d", createStatus)
	}
	if created == nil || created.Data.SubID == "" {
		t.Fatalf("Created transaction has empty SubID")
	}

	// List should return at least the transaction we created
	// Use large page_size to ensure our transaction appears even when
	// other tests have created many transactions before this test runs
	result, status, err := doJSONRequest[models.TransactionsResponse]("GET", "/transactions?page_size=1000", nil)
	if err != nil {
		t.Fatalf("Failed to list transactions: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	found := false
	for _, tx := range result.Data {
		if tx.SubID == created.Data.SubID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected created transaction %s to be in list", created.Data.SubID)
	}
}

func TestGetTransaction_ValidUUIDFormat_NotFound(t *testing.T) {
	_, status, _ := doGetTransaction("00000000-0000-0000-0000-000000000000")
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent UUID, got %d", status)
	}
}
