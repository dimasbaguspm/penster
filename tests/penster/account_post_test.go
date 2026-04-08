package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestCreateAccount_Success(t *testing.T) {
	req := &models.CreateAccountRequest{
		Name:    "Test Account",
		Type:    models.AccountTypeExpense,
		Balance: 1000,
	}
	result, status, err := doCreateAccount(req)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Name != "Test Account" {
		t.Errorf("Expected name 'Test Account', got %s", result.Data.Name)
	}
	if result.Data.Type != models.AccountTypeExpense {
		t.Errorf("Expected type 'expense', got %s", result.Data.Type)
	}
}

func TestCreateAccount_ValidationError_MissingName(t *testing.T) {
	req := &models.CreateAccountRequest{
		Type:    models.AccountTypeExpense,
		Balance: 1000,
	}
	_, status, _ := doCreateAccount(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateAccount_ValidationError_InvalidType(t *testing.T) {
	req := &models.CreateAccountRequest{
		Name:    "Test Account",
		Type:    "invalid_type",
		Balance: 1000,
	}
	result, status, err := doCreateAccount(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result != nil && result.Success {
		t.Errorf("Expected success=false for invalid type")
	}
}

func TestCreateAccount_ValidationError_EmptyName(t *testing.T) {
	req := &models.CreateAccountRequest{
		Name:    "",
		Type:    models.AccountTypeExpense,
		Balance: 1000,
	}
	result, status, err := doCreateAccount(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty name, got %d", status)
	}
	if result != nil && result.Success {
		t.Errorf("Expected success=false for empty name")
	}
}

func TestCreateAccount_ValidationError_NegativeBalance(t *testing.T) {
	req := &models.CreateAccountRequest{
		Name:    "Negative Balance Account",
		Type:    models.AccountTypeExpense,
		Balance: -100,
	}
	result, status, err := doCreateAccount(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for negative balance, got %d", status)
	}
	if result != nil && result.Success {
		t.Errorf("Expected success=false for negative balance")
	}
}

func TestCreateAccount_Success_IncomeType(t *testing.T) {
	req := &models.CreateAccountRequest{
		Name:    "Income Account",
		Type:    models.AccountTypeIncome,
		Balance: 5000,
	}
	result, status, err := doCreateAccount(req)
	if err != nil {
		t.Fatalf("Failed to create income account: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Type != models.AccountTypeIncome {
		t.Errorf("Expected type 'income', got %s", result.Data.Type)
	}
}

func TestCreateAccount_Success_TransferType(t *testing.T) {
	req := &models.CreateAccountRequest{
		Name:    "Transfer Account",
		Type:    models.AccountTypeTransfer,
		Balance: 0,
	}
	result, status, err := doCreateAccount(req)
	if err != nil {
		t.Fatalf("Failed to create transfer account: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Type != models.AccountTypeTransfer {
		t.Errorf("Expected type 'transfer', got %s", result.Data.Type)
	}
}
