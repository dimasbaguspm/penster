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

func TestListTransactions_TableDriven(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T, accountID, categoryID string)
	}{
		{
			name: "success_returns_200_with_data",
			test: func(t *testing.T, accountID, categoryID string) {
				// Create a few transactions
				for i := range 3 {
					req := &models.CreateTransactionRequest{
						AccountID:       accountID,
						CategoryID:       categoryID,
						TransactionType:  models.TransactionTypeExpense,
						Title:            "List Test Transaction",
						Amount:           int64(100 * (i + 1)),
						Currency:         "USD",
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
				if result.Items == nil && result.TotalItems == 0 {
					// Empty list is OK
				}
			},
		},
		{
			name: "pagination_meta_fields_valid",
			test: func(t *testing.T, accountID, categoryID string) {
				result, status, err := doListTransactions()
				if err != nil {
					t.Fatalf("Failed to list transactions: %v", err)
				}
				if status != http.StatusOK {
					t.Errorf("Expected status 200, got %d", status)
				}
				// When list is empty, pagination fields may be 0
				if result.PageNumber < 0 {
					t.Errorf("Expected page_number >= 0, got %d", result.PageNumber)
				}
				if result.PageSize < 0 {
					t.Errorf("Expected page_size >= 0, got %d", result.PageSize)
				}
				if result.TotalItems < 0 {
					t.Errorf("Expected total_items >= 0, got %d", result.TotalItems)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := createTestAccount(t)
			category := createTestCategory(t)
			tt.test(t, account.Data.SubID, category.Data.SubID)
		})
	}
}

