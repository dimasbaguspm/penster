package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestDeleteAccount_Success(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Account To Delete",
		Type:    models.AccountTypeExpense,
		Balance: 1000,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	result, status, err := doDeleteAccount(id)
	if err != nil {
		t.Fatalf("Failed to delete account: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}

	_, getStatus, _ := doGetAccount(id)
	if getStatus != http.StatusNotFound {
		t.Errorf("Expected deleted account to return 404 on GET, got %d", getStatus)
	}
}

func TestDeleteAccount_NotFound(t *testing.T) {
	nonExistentID := "00000000-0000-0000-0000-000000000000"
	_, status, _ := doDeleteAccount(nonExistentID)
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

func TestDeleteAccount_InvalidUUID(t *testing.T) {
	_, status, _ := doDeleteAccount("not-a-uuid")
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}

func TestDeleteAccount_AlreadyDeleted(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Double Delete Test",
		Type:    models.AccountTypeIncome,
		Balance: 5000,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	_, status1, _ := doDeleteAccount(id)
	if status1 != http.StatusOK {
		t.Errorf("Expected status 200 on first delete, got %d", status1)
	}

	_, status2, _ := doDeleteAccount(id)
	if status2 != http.StatusNotFound {
		t.Errorf("Expected status 404 on second delete, got %d", status2)
	}
}

func TestDeleteAccount_GetAfterDelete(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Get After Delete Test",
		Type:    models.AccountTypeTransfer,
		Balance: 0,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	_, _, _ = doDeleteAccount(id)

	_, status, _ := doGetAccount(id)
	if status != http.StatusNotFound {
		t.Errorf("Expected 404 after delete, got %d", status)
	}
}
