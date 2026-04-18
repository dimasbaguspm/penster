package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestDeleteTransaction_Success(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Transaction To Delete",
		Amount:          500,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	_, status, err := doDeleteTransaction(id)
	if err != nil {
		t.Fatalf("Failed to delete transaction: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	_, getStatus, _ := doGetTransaction(id)
	if getStatus != http.StatusNotFound {
		t.Errorf("Expected deleted transaction to return 404 on GET, got %d", getStatus)
	}
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	nonExistentID := "00000000-0000-0000-0000-000000000000"
	_, status, _ := doDeleteTransaction(nonExistentID)
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

func TestDeleteTransaction_InvalidUUID(t *testing.T) {
	_, status, _ := doDeleteTransaction("not-a-uuid")
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}

func TestDeleteTransaction_AlreadyDeleted(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeIncome,
		Title:           "Double Delete Test",
		Amount:          1000,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	_, status1, _ := doDeleteTransaction(id)
	if status1 != http.StatusOK {
		t.Errorf("Expected status 200 on first delete, got %d", status1)
	}

	_, status2, _ := doDeleteTransaction(id)
	if status2 != http.StatusNotFound {
		t.Errorf("Expected status 404 on second delete, got %d", status2)
	}
}

func TestDeleteTransaction_GetAfterDelete(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeTransfer,
		Title:           "Get After Delete Test",
		Amount:          250,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	_, _, _ = doDeleteTransaction(id)

	_, status, _ := doGetTransaction(id)
	if status != http.StatusNotFound {
		t.Errorf("Expected 404 after delete, got %d", status)
	}
}

func TestDeleteTransaction_ListAfterDelete(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "List After Delete Test",
		Amount:          300,
		Currency:        "EUR",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	_, _, _ = doDeleteTransaction(id)

	// List should not contain the deleted transaction
	result, status, err := doListTransactions()
	if err != nil {
		t.Fatalf("Failed to list transactions: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	for _, tx := range result.Items {
		if tx.SubID == id {
			t.Errorf("Deleted transaction should not appear in list")
		}
	}
}
