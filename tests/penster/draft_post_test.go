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
		TransactedAt:    "2024-01-15",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Failed to create draft: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
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
		TransactedAt:    "2024-01-16",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Failed to create income draft: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
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
		TransactedAt:      "2024-01-17",
		Source:            string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Failed to create transfer draft: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
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
		TransactedAt:    "2024-01-18",
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
		TransactedAt:    "2024-01-19",
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
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Notes != "This is a test note" {
		t.Errorf("Expected notes 'This is a test note', got %s", result.Data.Notes)
	}
}

// TestCreateDraft_ValidationError_MissingAccountID verifies validation when account_id is missing.
func TestCreateDraft_ValidationError_MissingAccountID(t *testing.T) {
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Missing Account",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-01-20",
		Source:          string(models.DraftSourceManual),
	}
	_, status, _ := doCreateDraft(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_ValidationError_MissingCategoryID verifies validation when category_id is missing.
func TestCreateDraft_ValidationError_MissingCategoryID(t *testing.T) {
	account := createTestDraftAccount(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Missing Category",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-01-21",
		Source:          string(models.DraftSourceManual),
	}
	_, status, _ := doCreateDraft(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_ValidationError_MissingTitle verifies validation when title is missing.
func TestCreateDraft_ValidationError_MissingTitle(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-01-22",
		Source:          string(models.DraftSourceManual),
	}
	_, status, _ := doCreateDraft(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_ValidationError_MissingAmount verifies validation when amount is missing.
func TestCreateDraft_ValidationError_MissingAmount(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Missing Amount",
		Currency:        "USD",
		TransactedAt:    "2024-01-23",
		Source:          string(models.DraftSourceManual),
	}
	_, status, _ := doCreateDraft(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_ValidationError_MissingCurrency verifies validation when currency is missing.
func TestCreateDraft_ValidationError_MissingCurrency(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Missing Currency",
		Amount:          100,
		TransactedAt:    "2024-01-24",
		Source:          string(models.DraftSourceManual),
	}
	_, status, _ := doCreateDraft(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_ValidationError_MissingTransactedAt verifies validation when transacted_at is missing.
func TestCreateDraft_ValidationError_MissingTransactedAt(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Missing TransactedAt",
		Amount:          100,
		Currency:        "USD",
		Source:          string(models.DraftSourceManual),
	}
	_, status, _ := doCreateDraft(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_ValidationError_MissingSource verifies validation when source is missing.
func TestCreateDraft_ValidationError_MissingSource(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Missing Source",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-01-25",
	}
	_, status, _ := doCreateDraft(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_ValidationError_ZeroAmount verifies validation when amount is zero.
func TestCreateDraft_ValidationError_ZeroAmount(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Zero Amount Draft",
		Amount:          0,
		Currency:        "USD",
		TransactedAt:    "2024-01-26",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false, got true")
	}
}

// TestCreateDraft_ValidationError_NegativeAmount verifies validation when amount is negative.
func TestCreateDraft_ValidationError_NegativeAmount(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Negative Amount Draft",
		Amount:          -100,
		Currency:        "USD",
		TransactedAt:    "2024-01-27",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false, got true")
	}
}

// TestCreateDraft_ValidationError_InvalidTransactionType verifies validation for invalid transaction type.
func TestCreateDraft_ValidationError_InvalidTransactionType(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: "invalid_type",
		Title:           "Invalid Type",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-01-28",
		Source:          string(models.DraftSourceManual),
	}
	status, err := doRequest("POST", "/drafts", req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_ValidationError_InvalidSource verifies validation for invalid source.
func TestCreateDraft_ValidationError_InvalidSource(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Invalid Source",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-01-29",
		Source:          "invalid_source",
	}
	status, err := doRequest("POST", "/drafts", req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_ValidationError_InvalidAccountID verifies validation for non-existent account.
func TestCreateDraft_ValidationError_InvalidAccountID(t *testing.T) {
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       "00000000-0000-0000-0000-000000000000",
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Invalid Account Draft",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-01-30",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false, got true")
	}
}

// TestCreateDraft_ValidationError_InvalidCategoryID verifies validation for non-existent category.
func TestCreateDraft_ValidationError_InvalidCategoryID(t *testing.T) {
	account := createTestDraftAccount(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      "00000000-0000-0000-0000-000000000000",
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Invalid Category Draft",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-01-31",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false, got true")
	}
}

// TestCreateDraft_ValidationError_InvalidTransferAccountID verifies validation for non-existent transfer account.
func TestCreateDraft_ValidationError_InvalidTransferAccountID(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:         account.Data.SubID,
		TransferAccountID: "00000000-0000-0000-0000-000000000000",
		CategoryID:        category.Data.SubID,
		TransactionType:   string(models.TransactionTypeTransfer),
		Title:             "Invalid Transfer Account Draft",
		Amount:            100,
		Currency:          "USD",
		TransactedAt:      "2024-02-01",
		Source:            string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false, got true")
	}
}

// TestCreateDraft_ValidationError_MalformedAccountUUID verifies validation for malformed account UUID.
func TestCreateDraft_ValidationError_MalformedAccountUUID(t *testing.T) {
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       "not-a-valid-uuid",
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Malformed UUID Draft",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-02-02",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false, got true")
	}
}

// TestCreateDraft_ValidationError_MalformedCategoryUUID verifies validation for malformed category UUID.
func TestCreateDraft_ValidationError_MalformedCategoryUUID(t *testing.T) {
	account := createTestDraftAccount(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      "not-a-valid-uuid",
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Malformed Category UUID Draft",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-02-03",
		Source:          string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false, got true")
	}
}

// TestCreateDraft_TransferToSameAccount verifies validation for transfer to same account.
func TestCreateDraft_TransferToSameAccount(t *testing.T) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:         account.Data.SubID,
		TransferAccountID: account.Data.SubID, // same account
		CategoryID:        category.Data.SubID,
		TransactionType:   string(models.TransactionTypeTransfer),
		Title:             "Transfer to Same Account",
		Amount:            100,
		Currency:          "USD",
		TransactedAt:      "2024-02-04",
		Source:            string(models.DraftSourceManual),
	}
	result, status, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false, got true")
	}
}

// TestCreateDraft_EmptyAccountID verifies validation for empty account_id.
func TestCreateDraft_EmptyAccountID(t *testing.T) {
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       "",
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Empty Account ID Draft",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-02-05",
		Source:          string(models.DraftSourceManual),
	}
	_, status, _ := doCreateDraft(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestCreateDraft_EmptyCategoryID verifies validation for empty category_id.
func TestCreateDraft_EmptyCategoryID(t *testing.T) {
	account := createTestDraftAccount(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      "",
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Empty Category ID Draft",
		Amount:          100,
		Currency:        "USD",
		TransactedAt:    "2024-02-06",
		Source:          string(models.DraftSourceManual),
	}
	_, status, _ := doCreateDraft(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}
