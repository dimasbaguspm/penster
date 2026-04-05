package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestConfirmDraft_Expense_DecreasesBalance verifies confirming an expense draft
// decreases the source account's balance.
func TestConfirmDraft_Expense_DecreasesBalance(t *testing.T) {
	account := createTestAccountWithBalance(t, 1000)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Balance Test Expense Draft",
		Amount:          300,
		Currency:        "USD",
		TransactedAt:    "2024-04-01",
		Source:          string(models.DraftSourceManual),
	}
	draft, _, _ := doCreateDraft(req)

	// Confirm the draft
	_, _, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm draft: %v", err)
	}

	// Verify account balance decreased
	updatedAccount, status, err := doGetAccount(account.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}
	expectedBalance := int64(1000) - 300
	if updatedAccount.Data.Balance != expectedBalance {
		t.Errorf("Expected balance %d after expense draft confirmed, got %d", expectedBalance, updatedAccount.Data.Balance)
	}
}

// TestConfirmDraft_Income_IncreasesBalance verifies confirming an income draft
// increases the source account's balance.
func TestConfirmDraft_Income_IncreasesBalance(t *testing.T) {
	account := createTestAccountWithBalance(t, 500)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeIncome),
		Title:           "Balance Test Income Draft",
		Amount:          200,
		Currency:        "USD",
		TransactedAt:    "2024-04-02",
		Source:          string(models.DraftSourceManual),
	}
	draft, _, _ := doCreateDraft(req)

	// Confirm the draft
	_, _, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm income draft: %v", err)
	}

	// Verify account balance increased
	updatedAccount, status, err := doGetAccount(account.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}
	expectedBalance := int64(500) + 200
	if updatedAccount.Data.Balance != expectedBalance {
		t.Errorf("Expected balance %d after income draft confirmed, got %d", expectedBalance, updatedAccount.Data.Balance)
	}
}

// TestConfirmDraft_Transfer_UpdatesBothAccounts verifies confirming a transfer draft
// decreases source account and increases destination account.
func TestConfirmDraft_Transfer_UpdatesBothAccounts(t *testing.T) {
	sourceAccount := createTestAccountWithBalance(t, 1000)
	destAccount := createTestAccountWithBalance(t, 500)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:         sourceAccount.Data.SubID,
		TransferAccountID: destAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   string(models.TransactionTypeTransfer),
		Title:             "Balance Test Transfer Draft",
		Amount:            250,
		Currency:          "USD",
		TransactedAt:      "2024-04-03",
		Source:            string(models.DraftSourceManual),
	}
	draft, _, _ := doCreateDraft(req)

	// Confirm the draft
	_, _, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm transfer draft: %v", err)
	}

	// Verify source account balance decreased
	updatedSource, status, err := doGetAccount(sourceAccount.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get source account: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}
	expectedSourceBalance := int64(1000) - 250
	if updatedSource.Data.Balance != expectedSourceBalance {
		t.Errorf("Expected source balance %d after transfer draft confirmed, got %d", expectedSourceBalance, updatedSource.Data.Balance)
	}

	// Verify destination account balance increased
	updatedDest, status, err := doGetAccount(destAccount.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get destination account: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}
	expectedDestBalance := int64(500) + 250
	if updatedDest.Data.Balance != expectedDestBalance {
		t.Errorf("Expected dest balance %d after transfer draft confirmed, got %d", expectedDestBalance, updatedDest.Data.Balance)
	}
}

// TestConfirmDraft_ThenDeleteTransaction_RestoresBalance verifies that
// confirming a draft creates a transaction, and if that transaction is deleted,
// the balance is restored.
func TestConfirmDraft_ThenDeleteTransaction_RestoresBalance(t *testing.T) {
	account := createTestAccountWithBalance(t, 1000)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Delete After Confirm Draft",
		Amount:          400,
		Currency:        "USD",
		TransactedAt:    "2024-04-04",
		Source:          string(models.DraftSourceManual),
	}
	draft, _, _ := doCreateDraft(req)

	// Confirm the draft
	txResult, _, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm draft: %v", err)
	}

	// Verify balance after confirm
	updatedAccount, _, _ := doGetAccount(account.Data.SubID)
	if updatedAccount.Data.Balance != 600 {
		t.Fatalf("Expected balance 600 after confirm, got %d", updatedAccount.Data.Balance)
	}

	// Delete the created transaction
	_, _, err = doDeleteTransaction(txResult.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to delete transaction: %v", err)
	}

	// Balance should be restored to original
	updatedAccount, _, _ = doGetAccount(account.Data.SubID)
	if updatedAccount.Data.Balance != 1000 {
		t.Errorf("Expected balance 1000 after deleting transaction, got %d", updatedAccount.Data.Balance)
	}
}

// TestConfirmDraft_Transfer_ThenDeleteTransaction_RestoresBothBalances verifies
// that confirming a transfer draft creates a transaction, and deleting that
// transaction restores both balances.
func TestConfirmDraft_Transfer_ThenDeleteTransaction_RestoresBothBalances(t *testing.T) {
	sourceAccount := createTestAccountWithBalance(t, 1000)
	destAccount := createTestAccountWithBalance(t, 500)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:         sourceAccount.Data.SubID,
		TransferAccountID: destAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   string(models.TransactionTypeTransfer),
		Title:             "Delete Transfer After Confirm Draft",
		Amount:            300,
		Currency:          "USD",
		TransactedAt:      "2024-04-05",
		Source:            string(models.DraftSourceManual),
	}
	draft, _, _ := doCreateDraft(req)

	// Confirm the draft
	txResult, _, err := doConfirmDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm transfer draft: %v", err)
	}

	// Verify balances after confirm
	sourceAfter, _, _ := doGetAccount(sourceAccount.Data.SubID)
	if sourceAfter.Data.Balance != 700 {
		t.Fatalf("Expected source balance 700 after confirm, got %d", sourceAfter.Data.Balance)
	}
	destAfter, _, _ := doGetAccount(destAccount.Data.SubID)
	if destAfter.Data.Balance != 800 {
		t.Fatalf("Expected dest balance 800 after confirm, got %d", destAfter.Data.Balance)
	}

	// Delete the created transaction
	_, _, err = doDeleteTransaction(txResult.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to delete transaction: %v", err)
	}

	// Both balances should be restored
	sourceAfter, _, _ = doGetAccount(sourceAccount.Data.SubID)
	if sourceAfter.Data.Balance != 1000 {
		t.Errorf("Expected source balance 1000 after deleting transaction, got %d", sourceAfter.Data.Balance)
	}
	destAfter, _, _ = doGetAccount(destAccount.Data.SubID)
	if destAfter.Data.Balance != 500 {
		t.Errorf("Expected dest balance 500 after deleting transaction, got %d", destAfter.Data.Balance)
	}
}

// TestRejectDraft_NoBalanceChange verifies that rejecting a draft does not affect balance.
func TestRejectDraft_NoBalanceChange(t *testing.T) {
	account := createTestAccountWithBalance(t, 1000)
	category := createTestDraftCategory(t)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Reject No Balance Change Draft",
		Amount:          500,
		Currency:        "USD",
		TransactedAt:    "2024-04-06",
		Source:          string(models.DraftSourceManual),
	}
	draft, _, _ := doCreateDraft(req)

	// Reject the draft
	_, _, err := doRejectDraft(draft.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to reject draft: %v", err)
	}

	// Verify account balance is unchanged
	updatedAccount, status, err := doGetAccount(account.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}
	if updatedAccount.Data.Balance != 1000 {
		t.Errorf("Expected balance 1000 (unchanged) after reject, got %d", updatedAccount.Data.Balance)
	}
}

// TestConfirmDraft_MultipleExpenses_AccumulatesBalance verifies confirming multiple
// expense drafts correctly accumulates balance changes.
func TestConfirmDraft_MultipleExpenses_AccumulatesBalance(t *testing.T) {
	account := createTestAccountWithBalance(t, 1000)
	category := createTestDraftCategory(t)

	// Create and confirm first expense draft
	req1 := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "First Expense Draft",
		Amount:          200,
		Currency:        "USD",
		TransactedAt:    "2024-04-07",
		Source:          string(models.DraftSourceManual),
	}
	draft1, _, _ := doCreateDraft(req1)
	_, _, err := doConfirmDraft(draft1.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm first draft: %v", err)
	}

	// Create and confirm second expense draft
	req2 := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           "Second Expense Draft",
		Amount:          300,
		Currency:        "USD",
		TransactedAt:    "2024-04-08",
		Source:          string(models.DraftSourceManual),
	}
	draft2, _, _ := doCreateDraft(req2)
	_, _, err = doConfirmDraft(draft2.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to confirm second draft: %v", err)
	}

	// Verify final balance is correct
	updatedAccount, _, _ := doGetAccount(account.Data.SubID)
	expectedBalance := int64(1000) - 200 - 300
	if updatedAccount.Data.Balance != expectedBalance {
		t.Errorf("Expected balance %d after two expense drafts confirmed, got %d", expectedBalance, updatedAccount.Data.Balance)
	}
}
