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
	if result.Data.Name != "Test Category" {
		t.Errorf("Expected name 'Test Category', got '%s'", result.Data.Name)
	}
	if result.Data.Type != models.CategoryTypeExpense {
		t.Errorf("Expected type 'expense', got '%s'", result.Data.Type)
	}
}

func TestCreateCategory_ValidationError(t *testing.T) {
	tests := []struct {
		name    string
		req     *models.CreateCategoryRequest
	}{
		{
			name: "Missing Name",
			req:  &models.CreateCategoryRequest{Type: models.CategoryTypeExpense},
		},
		{
			name: "Missing Type",
			req:  &models.CreateCategoryRequest{Name: "Missing Type Category"},
		},
		{
			name: "Invalid Type",
			req:  &models.CreateCategoryRequest{Name: "Invalid Type Category", Type: "invalid_type"},
		},
		{
			name: "Empty Name",
			req:  &models.CreateCategoryRequest{Name: "", Type: models.CategoryTypeExpense},
		},
		{
			name: "Whitespace Name",
			req:  &models.CreateCategoryRequest{Name: "   ", Type: models.CategoryTypeExpense},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, status, _ := doCreateCategory(tt.req)
			if status != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", status)
			}
		})
	}
}

func TestCreateCategory_Success_AllTypes(t *testing.T) {
	tests := []struct {
		name         string
		CategoryType models.CategoryType
	}{
		{name: "Expense", CategoryType: models.CategoryTypeExpense},
		{name: "Income", CategoryType: models.CategoryTypeIncome},
		{name: "Transfer", CategoryType: models.CategoryTypeTransfer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &models.CreateCategoryRequest{
				Name: tt.name + " Category",
				Type: tt.CategoryType,
			}
			result, status, err := doCreateCategory(req)
			if err != nil {
				t.Fatalf("Failed to create %s category: %v", tt.name, err)
			}
			if status != http.StatusCreated {
				t.Errorf("Expected status 201, got %d", status)
			}
			if result.Data.Type != tt.CategoryType {
				t.Errorf("Expected type '%s', got '%s'", tt.CategoryType, result.Data.Type)
			}
		})
	}
}
