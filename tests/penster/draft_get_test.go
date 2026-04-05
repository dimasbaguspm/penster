package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestGetDraft_Success verifies getting a single draft by ID.
func TestGetDraft_Success(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	result, status, err := doGetDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get draft: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.SubID != draft.Data.SubID {
		t.Errorf("Expected sub_id '%s', got %s", draft.Data.SubID, result.Data.SubID)
	}
	if result.Data.Title != draft.Data.Title {
		t.Errorf("Expected title '%s', got %s", draft.Data.Title, result.Data.Title)
	}
}

// TestGetDraft_NotFound verifies getting a non-existent draft returns 404.
func TestGetDraft_NotFound(t *testing.T) {
	_, status, _ := doGetDraft("00000000-0000-0000-0000-000000000000")
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

// TestGetDraft_MalformedUUID verifies getting a draft with malformed UUID returns 400.
func TestGetDraft_MalformedUUID(t *testing.T) {
	_, status, _ := doGetDraft("not-a-valid-uuid")
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestListDrafts_Success verifies listing all drafts.
func TestListDrafts_Success(t *testing.T) {
	// Create a few drafts
	_, _, _ = createTestDraftWithAccountAndCategory(t)
	_, _, _ = createTestDraftWithAccountAndCategory(t)

	result, status, err := doListDrafts()
	if err != nil {
		t.Fatalf("Failed to list drafts: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if len(result.Data) < 2 {
		t.Errorf("Expected at least 2 drafts, got %d", len(result.Data))
	}
}

// TestListDrafts_Empty verifies listing drafts when none exist.
func TestListDrafts_Empty(t *testing.T) {
	// Use a filter that won't match any drafts to ensure isolation
	// This test is flaky due to test isolation issues - other tests create drafts
	// that pollute the shared state
	result, status, err := doListDraftsFiltered("source", "nonexistent_source_filter")
	if err != nil {
		t.Fatalf("Failed to list drafts: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if len(result.Data) != 0 {
		t.Errorf("Expected 0 drafts with nonexistent source filter, got %d", len(result.Data))
	}
}

// TestListDrafts_FilterBySource verifies filtering drafts by source.
func TestListDrafts_FilterBySource(t *testing.T) {
	// Create manual draft
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)
	manualReq := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Manual Source Draft",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-02-01",
		Source:          string(models.DraftSourceManual),
	}
	doCreateDraft(manualReq)

	// Create ingestion draft
	ingestionReq := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Ingestion Source Draft",
		Amount:          200,
		Currency:        "USD",
		TransactedAt:    "2024-02-02",
		Source:          string(models.DraftSourceIngestion),
	}
	doCreateDraft(ingestionReq)

	// List only manual drafts
	result, status, err := doListDraftsFiltered("source", "manual")
	if err != nil {
		t.Fatalf("Failed to list drafts filtered by source: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	for _, d := range result.Data {
		if d.Source != string(models.DraftSourceManual) {
			t.Errorf("Expected source 'manual', got %s", d.Source)
		}
	}
}

// TestListDrafts_FilterByStatus verifies filtering drafts by status.
func TestListDrafts_FilterByStatus(t *testing.T) {
	// Create a draft (defaults to pending)
	_, _, _ = createTestDraftWithAccountAndCategory(t)

	// List only pending drafts
	result, status, err := doListDraftsFiltered("status", "pending")
	if err != nil {
		t.Fatalf("Failed to list drafts filtered by status: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	for _, d := range result.Data {
		if d.Status != string(models.DraftStatusPending) {
			t.Errorf("Expected status 'pending', got %s", d.Status)
		}
	}
}

// TestListDrafts_FilterBySourceAndStatus verifies filtering drafts by both source and status.
func TestListDrafts_FilterBySourceAndStatus(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	// Create manual pending draft
	manualReq := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Manual Pending Draft",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-02-03",
		Source:          string(models.DraftSourceManual),
	}
	doCreateDraft(manualReq)

	// List manual pending drafts
	result, status, err := doListDraftsFilteredMultiple("source", "manual", "status", "pending")
	if err != nil {
		t.Fatalf("Failed to list drafts filtered: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	for _, d := range result.Data {
		if d.Source != string(models.DraftSourceManual) {
			t.Errorf("Expected source 'manual', got %s", d.Source)
		}
		if d.Status != string(models.DraftStatusPending) {
			t.Errorf("Expected status 'pending', got %s", d.Status)
		}
	}
}

// doListDraftsFiltered returns drafts filtered by a single query parameter.
func doListDraftsFiltered(filterKey, filterValue string) (*models.DraftsResponse, int, error) {
	result, status, err := doJSONRequest[models.DraftsResponse]("GET", "/drafts?"+filterKey+"="+filterValue, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doListDraftsFilteredMultiple returns drafts filtered by multiple query parameters.
func doListDraftsFilteredMultiple(filterKey1, filterValue1, filterKey2, filterValue2 string) (*models.DraftsResponse, int, error) {
	result, status, err := doJSONRequest[models.DraftsResponse]("GET", "/drafts?"+filterKey1+"="+filterValue1+"&"+filterKey2+"="+filterValue2, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// TestListDrafts_PaginationMeta verifies pagination metadata in list response.
func TestListDrafts_PaginationMeta(t *testing.T) {
	_, _, _ = createTestDraftWithAccountAndCategory(t)

	result, status, err := doListDrafts()
	if err != nil {
		t.Fatalf("Failed to list drafts: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Meta == nil {
		t.Fatalf("Expected meta to be present")
	}
	if result.Meta.Page != 1 {
		t.Errorf("Expected page 1, got %d", result.Meta.Page)
	}
	if result.Meta.PerPage != 10 {
		t.Errorf("Expected per_page 10, got %d", result.Meta.PerPage)
	}
}
