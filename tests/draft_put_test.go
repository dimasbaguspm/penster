package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestUpdateDraft_Success verifies updating a draft's fields.
func TestUpdateDraft_Success(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	newTitle := "Updated Draft Title"
	newAmount := int64(750)
	req := &models.UpdateDraftRequest{
		Title:  &newTitle,
		Amount: &newAmount,
	}
	result, status, err := doUpdateDraft(draft.Data.SubID, req)
	if err != nil {
		t.Fatalf("Failed to update draft: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Title != newTitle {
		t.Errorf("Expected title '%s', got %s", newTitle, result.Data.Title)
	}
	if result.Data.Amount != newAmount {
		t.Errorf("Expected amount %d, got %d", newAmount, result.Data.Amount)
	}
}

// TestUpdateDraft_Success_ChangeAmount verifies changing just the amount.
func TestUpdateDraft_Success_ChangeAmount(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	newAmount := int64(300)
	req := &models.UpdateDraftRequest{
		Amount: &newAmount,
	}
	result, status, err := doUpdateDraft(draft.Data.SubID, req)
	if err != nil {
		t.Fatalf("Failed to update draft amount: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Amount != newAmount {
		t.Errorf("Expected amount %d, got %d", newAmount, result.Data.Amount)
	}
}

// TestUpdateDraft_Success_ChangeTitle verifies changing just the title.
func TestUpdateDraft_Success_ChangeTitle(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	newTitle := "New Title Only"
	req := &models.UpdateDraftRequest{
		Title: &newTitle,
	}
	result, status, err := doUpdateDraft(draft.Data.SubID, req)
	if err != nil {
		t.Fatalf("Failed to update draft title: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Title != newTitle {
		t.Errorf("Expected title '%s', got %s", newTitle, result.Data.Title)
	}
}

// TestUpdateDraft_Success_ChangeCurrency verifies changing the currency (re-fetches rate).
func TestUpdateDraft_Success_ChangeCurrency(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	newCurrency := "EUR"
	req := &models.UpdateDraftRequest{
		Currency: &newCurrency,
	}
	result, status, err := doUpdateDraft(draft.Data.SubID, req)
	if err != nil {
		t.Fatalf("Failed to update draft currency: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Currency != newCurrency {
		t.Errorf("Expected currency '%s', got %s", newCurrency, result.Data.Currency)
	}
}

// TestUpdateDraft_Success_ChangeNotes verifies changing notes.
func TestUpdateDraft_Success_ChangeNotes(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	newNotes := "Updated test notes"
	req := &models.UpdateDraftRequest{
		Notes: &newNotes,
	}
	result, status, err := doUpdateDraft(draft.Data.SubID, req)
	if err != nil {
		t.Fatalf("Failed to update draft notes: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Notes != newNotes {
		t.Errorf("Expected notes '%s', got %s", newNotes, result.Data.Notes)
	}
}

// TestUpdateDraft_ChangeAccount verifies changing the account_id.
func TestUpdateDraft_ChangeAccount(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)
	newAccount := createTestDraftAccount(t)

	req := &models.UpdateDraftRequest{
		AccountID: &newAccount.Data.SubID,
	}
	result, status, err := doUpdateDraft(draft.Data.SubID, req)
	if err != nil {
		t.Fatalf("Failed to update draft account: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.AccountID != newAccount.Data.SubID {
		t.Errorf("Expected account_id '%s', got %s", newAccount.Data.SubID, result.Data.AccountID)
	}
}

// TestUpdateDraft_ChangeCategory verifies changing the category_id.
func TestUpdateDraft_ChangeCategory(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)
	newCategory := createTestDraftCategory(t)

	req := &models.UpdateDraftRequest{
		CategoryID: &newCategory.Data.SubID,
	}
	result, status, err := doUpdateDraft(draft.Data.SubID, req)
	if err != nil {
		t.Fatalf("Failed to update draft category: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.CategoryID != newCategory.Data.SubID {
		t.Errorf("Expected category_id '%s', got %s", newCategory.Data.SubID, result.Data.CategoryID)
	}
}

// TestUpdateDraft_ChangeTransactionType verifies changing the transaction type.
func TestUpdateDraft_ChangeTransactionType(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	newType := string(models.TransactionTypeIncome)
	req := &models.UpdateDraftRequest{
		TransactionType: &newType,
	}
	result, status, err := doUpdateDraft(draft.Data.SubID, req)
	if err != nil {
		t.Fatalf("Failed to update draft transaction type: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.TransactionType != newType {
		t.Errorf("Expected transaction_type '%s', got %s", newType, result.Data.TransactionType)
	}
}

// TestUpdateDraft_NotFound verifies updating a non-existent draft returns 404.
func TestUpdateDraft_NotFound(t *testing.T) {
	newTitle := "Updated Title"
	req := &models.UpdateDraftRequest{
		Title: &newTitle,
	}
	_, status, _ := doUpdateDraft("00000000-0000-0000-0000-000000000000", req)
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

// TestUpdateDraft_MalformedUUID verifies updating with malformed UUID returns 400.
func TestUpdateDraft_MalformedUUID(t *testing.T) {
	newTitle := "Updated Title"
	req := &models.UpdateDraftRequest{
		Title: &newTitle,
	}
	_, status, _ := doUpdateDraft("not-a-valid-uuid", req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestUpdateDraft_InvalidAccountID verifies updating with invalid account_id returns 400.
func TestUpdateDraft_InvalidAccountID(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	invalidAccountID := "00000000-0000-0000-0000-000000000000"
	req := &models.UpdateDraftRequest{
		AccountID: &invalidAccountID,
	}
	_, status, _ := doUpdateDraft(draft.Data.SubID, req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestUpdateDraft_InvalidCategoryID verifies updating with invalid category_id returns 400.
func TestUpdateDraft_InvalidCategoryID(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	invalidCategoryID := "00000000-0000-0000-0000-000000000000"
	req := &models.UpdateDraftRequest{
		CategoryID: &invalidCategoryID,
	}
	_, status, _ := doUpdateDraft(draft.Data.SubID, req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestUpdateDraft_InvalidTransferAccountID verifies updating transfer account to invalid ID returns 400.
func TestUpdateDraft_InvalidTransferAccountID(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	invalidTransferAccountID := "00000000-0000-0000-0000-000000000000"
	req := &models.UpdateDraftRequest{
		TransferAccountID: &invalidTransferAccountID,
	}
	_, status, _ := doUpdateDraft(draft.Data.SubID, req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestUpdateDraft_TransferToSameAccount verifies updating transfer to same account returns 400.
func TestUpdateDraft_TransferToSameAccount(t *testing.T) {
	draft, account, _ := createTestDraftWithAccountAndCategory(t)

	txType := string(models.TransactionTypeTransfer)
	req := &models.UpdateDraftRequest{
		TransactionType:   &txType,
		TransferAccountID: &account.Data.SubID, // same as account_id
	}
	_, status, _ := doUpdateDraft(draft.Data.SubID, req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestUpdateDraft_ChangeToZeroAmount verifies changing amount to zero returns 400.
func TestUpdateDraft_ChangeToZeroAmount(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	zeroAmount := int64(0)
	req := &models.UpdateDraftRequest{
		Amount: &zeroAmount,
	}
	_, status, _ := doUpdateDraft(draft.Data.SubID, req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

// TestUpdateDraft_ChangeToNegativeAmount verifies changing amount to negative returns 400.
func TestUpdateDraft_ChangeToNegativeAmount(t *testing.T) {
	draft, _, _ := createTestDraftWithAccountAndCategory(t)

	negativeAmount := int64(-100)
	req := &models.UpdateDraftRequest{
		Amount: &negativeAmount,
	}
	_, status, _ := doUpdateDraft(draft.Data.SubID, req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}
