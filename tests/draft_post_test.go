package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestCreateDraft_Success_Expense verifies creating an expense draft.
func TestCreateDraft_Success_Expense(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Test Expense Draft",
		Amount:          500,
		Currency:        "USD",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Failed to create draft: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if result.Data.Title != "Test Expense Draft" {
		t.Errorf("Expected title 'Test Expense Draft', got %s", result.Data.Title)
	}
	if result.Data.TransactionType != string(models.TransactionTypeExpense) {
		t.Errorf("Expected type 'expense', got %s", result.Data.TransactionType)
	}
	if result.Data.Amount != 500 {
		t.Errorf("Expected amount 500, got %d", result.Data.Amount)
	}
	if result.Data.Status != string(models.DraftStatusPending) {
		t.Errorf("Expected status 'pending', got %s", result.Data.Status)
	}
	if result.Data.Source != string(models.DraftSourceManual) {
		t.Errorf("Expected source 'manual', got %s", result.Data.Source)
	}
}

// TestCreateDraft_Success_Income verifies creating an income draft.
func TestCreateDraft_Success_Income(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeIncome),
		Title:           "Test Income Draft",
		Amount:          1000,
		Currency:        "EUR",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Failed to create income draft: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if result.Data.TransactionType != string(models.TransactionTypeIncome) {
		t.Errorf("Expected type 'income', got %s", result.Data.TransactionType)
	}
}

// TestCreateDraft_Success_Transfer verifies creating a transfer draft.
func TestCreateDraft_Success_Transfer(t *testing.T) {
	account := createTestDraftAccount(t)
	targetAccount := createTestDraftTransferAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:         account.Data.SubID,
		TransferAccountID: targetAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   string(models.TransactionTypeTransfer),
		Title:             "Test Transfer Draft",
		Amount:            250,
		Currency:          "USD",
		Source:            string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Failed to create transfer draft: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if result.Data.TransactionType != string(models.TransactionTypeTransfer) {
		t.Errorf("Expected type 'transfer', got %s", result.Data.TransactionType)
	}
	if result.Data.TransferAccountID != targetAccount.Data.SubID {
		t.Errorf("Expected transfer_account_id '%s', got %s", targetAccount.Data.SubID, result.Data.TransferAccountID)
	}
}

// TestCreateDraft_Success_Ingestion verifies creating an ingestion source draft.
func TestCreateDraft_Success_Ingestion(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Ingested Draft",
		Amount:          300,
		Currency:        "USD",
		Source:          string(models.DraftSourceIngestion),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Failed to create ingestion draft: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if result.Data.Source != string(models.DraftSourceIngestion) {
		t.Errorf("Expected source 'ingestion', got %s", result.Data.Source)
	}
}

// TestCreateDraft_WithNotes verifies creating a draft with notes.
func TestCreateDraft_WithNotes(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Draft With Notes",
		Amount:          750,
		Currency:        "USD",
		Source:          string(models.DraftSourceManual),
		Notes:           "This is a test note",
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Failed to create draft with notes: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if result.Data.Notes != "This is a test note" {
		t.Errorf("Expected notes 'This is a test note', got %s", result.Data.Notes)
	}
}

// TestCreateDraft_ValidationError_TableDriven verifies validation errors using table-driven subtests.
func TestCreateDraft_ValidationError_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		req        *models.CreateDraftRequest
		wantStatus int
		wantErr    bool
	}{
		{
			name: "missing_account_id",
			req: &models.CreateDraftRequest{
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType:  string(models.TransactionTypeExpense),
				Title:           "Missing Account",
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing_category_id",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Missing Category",
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing_title",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing_amount",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Missing Amount",
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing_currency",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Missing Currency",
				Amount:          100,
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing_source",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Missing Source",
				Amount:          100,
				Currency:        "USD",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "zero_amount",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Zero Amount Draft",
				Amount:          0,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "negative_amount",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Negative Amount Draft",
				Amount:          -100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid_transaction_type",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: "invalid_type",
				Title:           "Invalid Type",
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid_source",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Invalid Source",
				Amount:          100,
				Currency:        "USD",
				Source:          "invalid_source",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid_account_id",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000000",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Invalid Account Draft",
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid_category_id",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "00000000-0000-0000-0000-000000000000",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Invalid Category Draft",
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid_transfer_account_id",
			req: &models.CreateDraftRequest{
				AccountID:         "00000000-0000-0000-0000-000000000001",
				TransferAccountID: "00000000-0000-0000-0000-000000000000",
				CategoryID:        "00000000-0000-0000-0000-000000000001",
				TransactionType:   string(models.TransactionTypeTransfer),
				Title:             "Invalid Transfer Account Draft",
				Amount:            100,
				Currency:          "USD",
				Source:            string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "malformed_account_uuid",
			req: &models.CreateDraftRequest{
				AccountID:       "not-a-valid-uuid",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Malformed UUID Draft",
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "malformed_category_uuid",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "not-a-valid-uuid",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Malformed Category UUID Draft",
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "transfer_to_same_account",
			req: &models.CreateDraftRequest{
				AccountID:         "00000000-0000-0000-0000-000000000001",
				TransferAccountID: "00000000-0000-0000-0000-000000000001",
				CategoryID:        "00000000-0000-0000-0000-000000000001",
				TransactionType:   string(models.TransactionTypeTransfer),
				Title:             "Transfer to Same Account",
				Amount:            100,
				Currency:          "USD",
				Source:            string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "empty_account_id",
			req: &models.CreateDraftRequest{
				AccountID:       "",
				CategoryID:      "00000000-0000-0000-0000-000000000001",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Empty Account ID Draft",
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "empty_category_id",
			req: &models.CreateDraftRequest{
				AccountID:       "00000000-0000-0000-0000-000000000001",
				CategoryID:      "",
				TransactionType: string(models.TransactionTypeExpense),
				Title:           "Empty Category ID Draft",
				Amount:          100,
				Currency:        "USD",
				Source:          string(models.DraftSourceManual),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, status, err := doCreateDraft(tt.req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			if status != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, status)
			}
		})
	}
}
