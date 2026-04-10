package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestConfirmDraft_Success_Expense verifies confirming an expense draft creates a transaction.
func TestConfirmDraft_Success_Expense(t *testing.T) {
	draft, account, category := createTestDraftWithAccountAndCategory(t)

	result, status, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm draft: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Title != draft.Data.Title {
		t.Errorf("Expected title '%s', got %s", draft.Data.Title, result.Data.Title)
	}
	if result.Data.TransactionType != models.TransactionTypeExpense {
		t.Errorf("Expected type 'expense', got %s", result.Data.TransactionType)
	}
	if result.Data.Amount != draft.Data.Amount {
		t.Errorf("Expected amount %d, got %d", draft.Data.Amount, result.Data.Amount)
	}
	if result.Data.AccountID != account.Data.SubID {
		t.Errorf("Expected account_id '%s', got %s", account.Data.SubID, result.Data.AccountID)
	}
	if result.Data.CategoryID != category.Data.SubID {
		t.Errorf("Expected category_id '%s', got %s", category.Data.SubID, result.Data.CategoryID)
	}
}

// TestConfirmDraft_Success_Income verifies confirming an income draft.
func TestConfirmDraft_Success_Income(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeIncome),
		Title:           "Income Draft",
		Amount:          1000,
		Currency:        "USD",
		Source:          string(models.DraftSourceManual),
	}
	draft, _, _ := doCreateDraft(req)

	result, status, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm income draft: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.TransactionType != models.TransactionTypeIncome {
		t.Errorf("Expected type 'income', got %s", result.Data.TransactionType)
	}
}

// TestConfirmDraft_Success_Transfer verifies confirming a transfer draft.
func TestConfirmDraft_Success_Transfer(t *testing.T) {
	account := createTestDraftAccount(t)
	targetAccount := createTestDraftTransferAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:         account.Data.SubID,
		TransferAccountID: targetAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   string(models.TransactionTypeTransfer),
		Title:             "Transfer Draft",
		Amount:            250,
		Currency:          "USD",
		Source:            string(models.DraftSourceManual),
	}
	draft, _, _ := doCreateDraft(req)

	result, status, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm transfer draft: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.TransactionType != models.TransactionTypeTransfer {
		t.Errorf("Expected type 'transfer', got %s", result.Data.TransactionType)
	}
	if result.Data.TransferAccountID != targetAccount.Data.SubID {
		t.Errorf("Expected transfer_account_id '%s', got %s", targetAccount.Data.SubID, result.Data.TransferAccountID)
	}
}

// TestConfirmDraft_StatusUpdate_UpdatesDraftStatus verifies confirming updates draft status to confirmed.
func TestConfirmDraft_StatusUpdate_UpdatesDraftStatus(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	_, _, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm draft: %v", err)
	}

	// Verify draft status is now confirmed
	updatedDraft, status, err := doGetDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get draft after confirm: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if updatedDraft.Data.Status != string(models.DraftStatusConfirmed) {
		t.Errorf("Expected draft status 'confirmed', got '%s'", updatedDraft.Data.Status)
	}
}

// TestConfirmDraft_Error_TableDriven verifies error cases using table-driven subtests.
func TestConfirmDraft_Error_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() string
		draftID    string
		wantStatus int
	}{
		{
			name:       "not_found",
			setup:      func() string { return "00000000-0000-0000-0000-000000000000" },
			draftID:    "",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "malformed_uuid",
			setup:      func() string { return "not-a-valid-uuid" },
			draftID:    "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "already_confirmed",
			setup: func() string {
				draft, _, _ := createTestDraftWithAccountAndCategory(t)
				doConfirmDraft(draft.Data.SubID)
				return draft.Data.SubID
			},
			draftID:    "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "already_rejected",
			setup: func() string {
				draft, _, _ := createTestDraftWithAccountAndCategory(t)
				doRejectDraft(draft.Data.SubID)
				return draft.Data.SubID
			},
			draftID:    "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			draftID := tt.setup()
			result, status, _ := doConfirmDraft(draftID)
			if status != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, status)
			}
			if result != nil && result.Success {
				t.Errorf("Expected success=false, got true")
			}
		})
	}
}
