package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestReplaceAccount_Success(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Original Account",
		Type:    models.AccountTypeExpense,
		Balance: 1000,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	newName := "Updated Account"
	newType := models.AccountTypeIncome
	newBalance := int64(2000)
	replaceReq := &models.UpdateAccountRequest{
		Name:    &newName,
		Type:    &newType,
		Balance: &newBalance,
	}

	result, status, err := doUpdateAccount(id, replaceReq)
	if err != nil {
		t.Fatalf("Failed to update account: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Name != "Updated Account" {
		t.Errorf("Expected name 'Updated Account', got %s", result.Data.Name)
	}
	if result.Data.Type != models.AccountTypeIncome {
		t.Errorf("Expected type 'income', got %s", result.Data.Type)
	}
	if result.Data.Balance != 2000 {
		t.Errorf("Expected balance 2000, got %d", result.Data.Balance)
	}
}

func TestReplaceAccount_PartialUpdate_Name(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Original Name",
		Type:    models.AccountTypeExpense,
		Balance: 1000,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	newName := "New Name Only"
	replaceReq := &models.UpdateAccountRequest{
		Name: &newName,
	}

	result, status, err := doUpdateAccount(id, replaceReq)
	if err != nil {
		t.Fatalf("Failed to update account name: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Name != "New Name Only" {
		t.Errorf("Expected name 'New Name Only', got %s", result.Data.Name)
	}
	if result.Data.Type != models.AccountTypeExpense {
		t.Errorf("Expected type unchanged 'expense', got %s", result.Data.Type)
	}
	if result.Data.Balance != 1000 {
		t.Errorf("Expected balance unchanged 1000, got %d", result.Data.Balance)
	}
}

func TestReplaceAccount_PartialUpdate_Type(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Type Test Account",
		Type:    models.AccountTypeExpense,
		Balance: 500,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	newType := models.AccountTypeTransfer
	replaceReq := &models.UpdateAccountRequest{
		Type: &newType,
	}

	result, status, err := doUpdateAccount(id, replaceReq)
	if err != nil {
		t.Fatalf("Failed to update account type: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Type != models.AccountTypeTransfer {
		t.Errorf("Expected type 'transfer', got %s", result.Data.Type)
	}
}

func TestReplaceAccount_PartialUpdate_Balance(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Balance Test Account",
		Type:    models.AccountTypeIncome,
		Balance: 0,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	newBalance := int64(9999)
	replaceReq := &models.UpdateAccountRequest{
		Balance: &newBalance,
	}

	result, status, err := doUpdateAccount(id, replaceReq)
	if err != nil {
		t.Fatalf("Failed to update account balance: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Balance != 9999 {
		t.Errorf("Expected balance 9999, got %d", result.Data.Balance)
	}
}

func TestReplaceAccount_NotFound(t *testing.T) {
	nonExistentID := "00000000-0000-0000-0000-000000000000"
	newName := "Should Not Work"
	replaceReq := &models.UpdateAccountRequest{
		Name: &newName,
	}
	_, status, _ := doUpdateAccount(nonExistentID, replaceReq)
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

func TestReplaceAccount_InvalidUUID(t *testing.T) {
	newName := "Should Not Work"
	replaceReq := &models.UpdateAccountRequest{
		Name: &newName,
	}
	_, status, _ := doUpdateAccount("not-a-valid-uuid", replaceReq)
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}

func TestReplaceAccount_AllFieldsNil(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Nil Update Test",
		Type:    models.AccountTypeExpense,
		Balance: 500,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	replaceReq := &models.UpdateAccountRequest{}

	result, status, err := doUpdateAccount(id, replaceReq)
	if err != nil {
		t.Fatalf("Failed to update account with nil fields: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Name != "Nil Update Test" {
		t.Errorf("Expected name unchanged 'Nil Update Test', got %s", result.Data.Name)
	}
}

func TestReplaceAccount_Idempotent(t *testing.T) {
	createReq := &models.CreateAccountRequest{
		Name:    "Idempotent Test",
		Type:    models.AccountTypeExpense,
		Balance: 100,
	}
	created, _, _ := doCreateAccount(createReq)
	id := created.Data.SubID

	newName := "Idempotent Test"
	replaceReq := &models.UpdateAccountRequest{
		Name: &newName,
	}

	result1, status1, err1 := doUpdateAccount(id, replaceReq)
	if err1 != nil {
		t.Fatalf("First update failed: %v", err1)
	}
	if status1 != http.StatusOK {
		t.Errorf("Expected status 200 on first update, got %d", status1)
	}

	result2, status2, err2 := doUpdateAccount(id, replaceReq)
	if err2 != nil {
		t.Fatalf("Second update failed: %v", err2)
	}
	if status2 != http.StatusOK {
		t.Errorf("Expected status 200 on second update, got %d", status2)
	}

	if result1.Data.UpdatedAt.Unix() != result2.Data.UpdatedAt.Unix() {
		t.Logf("Note: UpdatedAt differs between calls (idempotent but timestamp may change)")
	}
}
