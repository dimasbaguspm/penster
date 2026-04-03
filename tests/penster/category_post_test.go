package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestCreateCategory_Success(t *testing.T) {
	req := &models.CreateCategoryRequest{
		Name: "Test Category",
		Type: models.CategoryTypeExpense,
	}
	result, status, err := doCreateCategory(req)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Name != "Test Category" {
		t.Errorf("Expected name 'Test Category', got '%s'", result.Data.Name)
	}
	if result.Data.Type != models.CategoryTypeExpense {
		t.Errorf("Expected type 'expense', got '%s'", result.Data.Type)
	}
}

func TestCreateCategory_ValidationError_MissingName(t *testing.T) {
	req := &models.CreateCategoryRequest{
		Type: models.CategoryTypeExpense,
	}
	_, status, _ := doCreateCategory(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateCategory_ValidationError_MissingType(t *testing.T) {
	req := &models.CreateCategoryRequest{
		Name: "Missing Type Category",
	}
	_, status, _ := doCreateCategory(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateCategory_ValidationError_InvalidType(t *testing.T) {
	req := &models.CreateCategoryRequest{
		Name: "Invalid Type Category",
		Type: "invalid_type",
	}
	status, err := doRequest("POST", "/categories", req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateCategory_Success_ExpenseType(t *testing.T) {
	req := &models.CreateCategoryRequest{
		Name: "Expense Category",
		Type: models.CategoryTypeExpense,
	}
	result, status, err := doCreateCategory(req)
	if err != nil {
		t.Fatalf("Failed to create expense category: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Type != models.CategoryTypeExpense {
		t.Errorf("Expected type 'expense', got '%s'", result.Data.Type)
	}
}

func TestCreateCategory_Success_IncomeType(t *testing.T) {
	req := &models.CreateCategoryRequest{
		Name: "Income Category",
		Type: models.CategoryTypeIncome,
	}
	result, status, err := doCreateCategory(req)
	if err != nil {
		t.Fatalf("Failed to create income category: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Type != models.CategoryTypeIncome {
		t.Errorf("Expected type 'income', got '%s'", result.Data.Type)
	}
}

func TestCreateCategory_Success_TransferType(t *testing.T) {
	req := &models.CreateCategoryRequest{
		Name: "Transfer Category",
		Type: models.CategoryTypeTransfer,
	}
	result, status, err := doCreateCategory(req)
	if err != nil {
		t.Fatalf("Failed to create transfer category: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Type != models.CategoryTypeTransfer {
		t.Errorf("Expected type 'transfer', got '%s'", result.Data.Type)
	}
}
