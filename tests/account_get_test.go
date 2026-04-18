package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestGetAccount_Success(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Test Account",
		Type:    models.AccountTypeExpense,
		Balance: 1000,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	result, status, err := doGetAccount(id)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.SubID != id {
		t.Errorf("Expected sub_id '%s', got '%s'", id, result.Data.SubID)
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	_, status, _ := doGetAccount("00000000-0000-0000-0000-000000000000")
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

func TestListAccounts_Success(t *testing.T) {
	result, status, err := doListAccounts()
	if err != nil {
		t.Fatalf("Failed to list accounts: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Items == nil && len(result.Items) == 0 {
		// Empty list is OK, Items is nil or empty
	}
}

func TestGetAccount_InvalidUUID(t *testing.T) {
	_, status, _ := doGetAccount("not-a-uuid")
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}

func TestListAccounts_PaginationMeta(t *testing.T) {
	result, status, err := doListAccounts()
	if err != nil {
		t.Fatalf("Failed to list accounts: %v", err)
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
}
