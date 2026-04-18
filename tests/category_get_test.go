package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestListCategories_Success(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "List Test Category",
		Type: models.CategoryTypeExpense,
	}
	_, _, _ = doCreateCategory(createReq)

	result, status, err := doListCategories()
	if err != nil {
		t.Fatalf("Failed to list categories: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	// Items may be nil or empty when list is empty, but should have data when items exist
	if result.Items == nil && result.TotalItems > 0 {
		t.Errorf("Expected items to be non-nil when totalItems=%d", result.TotalItems)
	}
}

func TestListCategories_Empty(t *testing.T) {
	_, status, err := doListCategories()
	if err != nil {
		t.Fatalf("Failed to list categories: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
}

func TestListCategories_PaginationMeta(t *testing.T) {
	result, status, err := doListCategories()
	if err != nil {
		t.Fatalf("Failed to list categories: %v", err)
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

func TestGetCategory_Success(t *testing.T) {
	createReq := &models.CreateCategoryRequest{
		Name: "Get Test Category",
		Type: models.CategoryTypeIncome,
	}
	created, _, _ := doCreateCategory(createReq)
	id := created.Data.SubID

	result, status, err := doGetCategory(id)
	if err != nil {
		t.Fatalf("Failed to get category: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Name != createReq.Name {
		t.Errorf("Expected name '%s', got '%s'", createReq.Name, result.Data.Name)
	}
	if result.Data.Type != createReq.Type {
		t.Errorf("Expected type '%s', got '%s'", createReq.Type, result.Data.Type)
	}
}

func TestGetCategory_NotFound(t *testing.T) {
	_, status, _ := doGetCategory("00000000-0000-0000-0000-000000000000")
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

func TestGetCategory_InvalidUUID(t *testing.T) {
	_, status, _ := doGetCategory("not-a-uuid")
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}
