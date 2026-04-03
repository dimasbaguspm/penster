package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestDeleteCategory_Success(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Category To Delete",
		Type: models.CategoryTypeExpense,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	result, status, err := doDeleteCategory(id)
	if err != nil {
		t.Fatalf("Failed to delete category: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}

	_, getStatus, _ := doGetCategory(id)
	if getStatus != http.StatusNotFound {
		t.Errorf("Expected deleted category to return 404 on GET, got %d", getStatus)
	}
}

func TestDeleteCategory_NotFound(t *testing.T) {
	nonExistentID := "00000000-0000-0000-0000-000000000000"
	_, status, _ := doDeleteCategory(nonExistentID)
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

func TestDeleteCategory_InvalidUUID(t *testing.T) {
	_, status, _ := doDeleteCategory("not-a-uuid")
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}

func TestDeleteCategory_AlreadyDeleted(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Double Delete Test",
		Type: models.CategoryTypeIncome,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	_, status1, _ := doDeleteCategory(id)
	if status1 != http.StatusOK {
		t.Errorf("Expected status 200 on first delete, got %d", status1)
	}

	_, status2, _ := doDeleteCategory(id)
	if status2 != http.StatusNotFound {
		t.Errorf("Expected status 404 on second delete, got %d", status2)
	}
}

func TestDeleteCategory_GetAfterDelete(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Get After Delete Test",
		Type: models.CategoryTypeTransfer,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	_, _, _ = doDeleteCategory(id)

	_, status, _ := doGetCategory(id)
	if status != http.StatusNotFound {
		t.Errorf("Expected 404 after delete, got %d", status)
	}
}
