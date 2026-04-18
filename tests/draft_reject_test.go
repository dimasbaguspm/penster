package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestRejectDraft_Success verifies rejecting a pending draft.
func TestRejectDraft_Success(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	_, status, err := doRejectDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to reject draft: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	// Verify draft status is now rejected
	updatedDraft, status, err := doGetDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get draft after reject: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if updatedDraft.Data.Status != string(models.DraftStatusRejected) {
		t.Errorf("Expected draft status 'rejected', got '%s'", updatedDraft.Data.Status)
	}
}

// TestRejectDraft_NotFound verifies rejecting a non-existent draft returns 404.
func TestRejectDraft_NotFound(t *testing.T) {
	_, status, _ := doRejectDraft("00000000-0000-0000-0000-000000000000")
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

// TestRejectDraft_MalformedUUID verifies rejecting with malformed UUID returns 400.
func TestRejectDraft_MalformedUUID(t *testing.T) {
	_, status, _ := doRejectDraft("not-a-valid-uuid")
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestRejectDraft_AlreadyRejected verifies rejecting an already rejected draft returns 400.
func TestRejectDraft_AlreadyRejected(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	// First reject
	_, _, err := doRejectDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to reject draft first time: %v", err)
	}

	// Second reject should fail
	_, status, _ := doRejectDraft(draft.Data.SubID)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestRejectDraft_AlreadyConfirmed verifies rejecting an already confirmed draft returns 400.
func TestRejectDraft_AlreadyConfirmed(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	// First confirm
	_, _, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm draft: %v", err)
	}

	// Reject should fail
	_, status, _ := doRejectDraft(draft.Data.SubID)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestRejectDraft_DoesNotCreateTransaction verifies rejecting a draft does not create a transaction.
func TestRejectDraft_DoesNotCreateTransaction(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	_, _, err := doRejectDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to reject draft: %v", err)
	}

	// List transactions - should not include the draft's data as a transaction
	txResult, _, err := doListTransactions()
	if err != nil {
		t.Fatalf("Failed to list transactions: %v", err)
	}

	// Find if any transaction has the same title as the draft
	for _, tx := range txResult.Items {
		if tx.Title == draft.Data.Title {
			t.Errorf("Draft was rejected but a transaction was created with the same title")
		}
	}
}
