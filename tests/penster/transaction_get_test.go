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
				if !result.Success {
					t.Errorf("Expected success=true, got false with error: %s", result.Error)
				}
				if result.Data == nil {
					t.Errorf("Expected data to be non-nil")
				}
				if result.Meta == nil {
					t.Errorf("Expected meta to be non-nil")
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

