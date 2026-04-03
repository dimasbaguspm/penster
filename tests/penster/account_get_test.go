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
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
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

func TestGetAccount_InvalidUUID(t *testing.T) {
	_, status, _ := doGetAccount("not-a-uuid")
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}

func TestGetAccount_ValidUUIDFormat_NotFound(t *testing.T) {
	_, status, _ := doGetAccount("00000000-0000-0000-0000-000000000000")
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent UUID, got %d", status)
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
