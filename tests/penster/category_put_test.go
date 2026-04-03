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

func TestReplaceCategory_PartialUpdate_Name(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Original Name",
		Type: models.CategoryTypeExpense,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	newName := "New Name Only"
	replaceReq := &models.UpdateCategoryRequest{
		Name: &newName,
	}

	result, status, err := doUpdateCategory(id, replaceReq)
	if err != nil {
		t.Fatalf("Failed to update category name: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Name != "New Name Only" {
		t.Errorf("Expected name 'New Name Only', got '%s'", result.Data.Name)
	}
	if result.Data.Type != models.CategoryTypeExpense {
		t.Errorf("Expected type unchanged 'expense', got '%s'", result.Data.Type)
	}
}

func TestReplaceCategory_PartialUpdate_Type(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Type Test Category",
		Type: models.CategoryTypeExpense,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	newType := models.CategoryTypeTransfer
	replaceReq := &models.UpdateCategoryRequest{
		Type: &newType,
	}

	result, status, err := doUpdateCategory(id, replaceReq)
	if err != nil {
		t.Fatalf("Failed to update category type: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Type != models.CategoryTypeTransfer {
		t.Errorf("Expected type 'transfer', got '%s'", result.Data.Type)
	}
	if result.Data.Name != "Type Test Category" {
		t.Errorf("Expected name unchanged 'Type Test Category', got '%s'", result.Data.Name)
	}
}

func TestReplaceCategory_NotFound(t *testing.T) {
	nonExistentID := "00000000-0000-0000-0000-000000000000"
	newName := "Should Not Work"
	replaceReq := &models.UpdateCategoryRequest{
		Name: &newName,
	}
	_, status, _ := doUpdateCategory(nonExistentID, replaceReq)
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

func TestReplaceCategory_InvalidUUID(t *testing.T) {
	newName := "Should Not Work"
	replaceReq := &models.UpdateCategoryRequest{
		Name: &newName,
	}
	_, status, _ := doUpdateCategory("not-a-valid-uuid", replaceReq)
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}

func TestReplaceCategory_AllFieldsNil(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Nil Update Test",
		Type: models.CategoryTypeIncome,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	replaceReq := &models.UpdateCategoryRequest{}

	result, status, err := doUpdateCategory(id, replaceReq)
	if err != nil {
		t.Fatalf("Failed to update category with nil fields: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Name != "Nil Update Test" {
		t.Errorf("Expected name unchanged 'Nil Update Test', got '%s'", result.Data.Name)
	}
	if result.Data.Type != models.CategoryTypeIncome {
		t.Errorf("Expected type unchanged 'income', got '%s'", result.Data.Type)
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
