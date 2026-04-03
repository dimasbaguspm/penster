package main

import (
	"fmt"
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
		fmt.Println("er", err)
		t.Fatalf("Failed to list categories: %v", err)
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
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
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

func TestGetCategory_ValidUUIDFormat_NotFound(t *testing.T) {
	_, status, _ := doGetCategory("00000000-0000-0000-0000-000000000000")
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent UUID, got %d", status)
	}
}
