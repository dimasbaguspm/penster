package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestReplaceCategory_Success(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Original Category",
		Type: models.CategoryTypeExpense,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	newName := "Updated Category"
	newType := models.CategoryTypeIncome
	replaceReq := &models.UpdateCategoryRequest{
		Name: &newName,
		Type: &newType,
	}

	result, status, err := doUpdateCategory(id, replaceReq)
	if err != nil {
		t.Fatalf("Failed to update category: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Name != "Updated Category" {
		t.Errorf("Expected name 'Updated Category', got '%s'", result.Data.Name)
	}
	if result.Data.Type != models.CategoryTypeIncome {
		t.Errorf("Expected type 'income', got '%s'", result.Data.Type)
	}
}

func TestReplaceCategory_PartialUpdate(t *testing.T) {
	tests := []struct {
		name          string
		newName       *string
		newType       *models.CategoryType
		expectedName  string
		expectedType  models.CategoryType
		unchangedName string
		unchangedType models.CategoryType
	}{
		{
			name:          "Update Name Only",
			newName:       strPtr("New Name Only"),
			expectedName:  "New Name Only",
			unchangedName: "Partial Update Test",
			unchangedType: models.CategoryTypeExpense,
		},
		{
			name:          "Update Type Only",
			newType:       ptr(models.CategoryTypeTransfer),
			unchangedName: "Partial Update Test",
			unchangedType: models.CategoryTypeTransfer,
		},
		{
			name:          "All Fields Nil",
			unchangedName: "Nil Update Test",
			unchangedType: models.CategoryTypeIncome,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createReq := &models.CreateCategoryRequest{
				Name: tt.unchangedName,
				Type: tt.unchangedType,
			}
			created, _, _ := doCreateCategory(createReq)
			id := created.Data.SubID

			replaceReq := &models.UpdateCategoryRequest{}
			if tt.newName != nil {
				replaceReq.Name = tt.newName
			}
			if tt.newType != nil {
				replaceReq.Type = tt.newType
			}

			result, status, err := doUpdateCategory(id, replaceReq)
			if err != nil {
				t.Fatalf("Failed to update category: %v", err)
			}
			if status != http.StatusOK {
				t.Errorf("Expected status 200, got %d", status)
			}

			expectedName := tt.expectedName
			if expectedName == "" {
				expectedName = tt.unchangedName
			}
			if result.Data.Name != expectedName {
				t.Errorf("Expected name '%s', got '%s'", expectedName, result.Data.Name)
			}
			if result.Data.Type != tt.unchangedType {
				t.Errorf("Expected type '%s', got '%s'", tt.unchangedType, result.Data.Type)
			}
		})
	}
}

func TestReplaceCategory_NotFoundOrInvalid(t *testing.T) {
	tests := []struct {
		name       string
		id         string
	}{
		{name: "Not Found", id: "00000000-0000-0000-0000-000000000000"},
		{name: "Invalid UUID", id: "not-a-valid-uuid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newName := "Should Not Work"
			replaceReq := &models.UpdateCategoryRequest{Name: &newName}
			_, status, _ := doUpdateCategory(tt.id, replaceReq)
			if status != http.StatusBadRequest && status != http.StatusNotFound {
				t.Errorf("Expected status 400 or 404, got %d", status)
			}
		})
	}
}

func TestReplaceCategory_ValidationError(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Valid Name",
		Type: models.CategoryTypeExpense,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	tests := []struct {
		name string
		req  *models.UpdateCategoryRequest
	}{
		{
			name: "Empty Name",
			req:  &models.UpdateCategoryRequest{Name: strPtr("")},
		},
		{
			name: "Whitespace Name",
			req:  &models.UpdateCategoryRequest{Name: strPtr("   ")},
		},
		{
			name: "Invalid Type",
			req:  &models.UpdateCategoryRequest{Type: ptr(models.CategoryType("invalid_type"))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, status, err := doUpdateCategory(id, tt.req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			if status != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", status)
			}
			if result != nil && result.Success {
				t.Errorf("Expected success=false")
			}
		})
	}
}

func TestReplaceCategory_Idempotent(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Idempotent Test",
		Type: models.CategoryTypeExpense,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	newName := "Idempotent Test"
	replaceReq := &models.UpdateCategoryRequest{
		Name: &newName,
	}

	result1, status1, err1 := doUpdateCategory(id, replaceReq)
	if err1 != nil {
		t.Fatalf("First update failed: %v", err1)
	}
	if status1 != http.StatusOK {
		t.Errorf("Expected status 200 on first update, got %d", status1)
	}

	result2, status2, err2 := doUpdateCategory(id, replaceReq)
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
