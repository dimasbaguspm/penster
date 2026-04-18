package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// TestCreateExpense_DecreasesBalance verifies that creating an expense
// transaction decreases the source account's balance.
func TestCreateExpense_DecreasesBalance(t *testing.T) {
	account := createTestAccountWithBalance(t, 1000)
	category := createTestCategory(t)

	expenseAmount := int64(300)
	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Balance Test Expense",
		Amount:          expenseAmount,
		Currency:        "USD",
	}
	_, _, err := doCreateTransaction(req)
	if err != nil {
		t.Fatalf("Failed to create expense transaction: %v", err)
	}

	// Verify account balance decreased
	updatedAccount, status, err := doGetAccount(account.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}
	expectedBalance := int64(1000) - expenseAmount
	if updatedAccount.Data.Balance != expectedBalance {
		t.Errorf("Expected balance %d after expense, got %d", expectedBalance, updatedAccount.Data.Balance)
	}
}

// TestCreateIncome_IncreasesBalance verifies that creating an income
// transaction increases the source account's balance.
func TestCreateIncome_IncreasesBalance(t *testing.T) {
	account := createTestAccountWithBalance(t, 500)
	category := createTestCategory(t)

	incomeAmount := int64(200)
	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeIncome,
		Title:           "Balance Test Income",
		Amount:          incomeAmount,
		Currency:        "USD",
	}
	_, _, err := doCreateTransaction(req)
	if err != nil {
		t.Fatalf("Failed to create income transaction: %v", err)
	}

	// Verify account balance increased
	updatedAccount, status, err := doGetAccount(account.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}
	expectedBalance := int64(500) + incomeAmount
	if updatedAccount.Data.Balance != expectedBalance {
		t.Errorf("Expected balance %d after income, got %d", expectedBalance, updatedAccount.Data.Balance)
	}
}

// TestCreateTransfer_UpdatesBothAccounts verifies that creating a transfer
// decreases source account and increases destination account.
func TestCreateTransfer_UpdatesBothAccounts(t *testing.T) {
	sourceAccount := createTestAccountWithBalance(t, 1000)
	destAccount := createTestAccountWithBalance(t, 500)
	category := createTestCategory(t)

	transferAmount := int64(250)
	req := &models.CreateTransactionRequest{
		AccountID:         sourceAccount.Data.SubID,
		TransferAccountID: destAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   models.TransactionTypeTransfer,
		Title:             "Balance Test Transfer",
		Amount:            transferAmount,
		Currency:          "USD",
	}
	_, _, err := doCreateTransaction(req)
	if err != nil {
		t.Fatalf("Failed to create transfer transaction: %v", err)
	}

	// Verify source account balance decreased
	updatedSource, status, err := doGetAccount(sourceAccount.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get source account: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}
	expectedSourceBalance := int64(1000) - transferAmount
	if updatedSource.Data.Balance != expectedSourceBalance {
		t.Errorf("Expected source balance %d after transfer, got %d", expectedSourceBalance, updatedSource.Data.Balance)
	}

	// Verify destination account balance increased
	updatedDest, status, err := doGetAccount(destAccount.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to get destination account: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}
	expectedDestBalance := int64(500) + transferAmount
	if updatedDest.Data.Balance != expectedDestBalance {
		t.Errorf("Expected dest balance %d after transfer, got %d", expectedDestBalance, updatedDest.Data.Balance)
	}
}

// TestUpdateTransaction_AmountChange_UpdatesBalance verifies that changing
// the amount of a transaction correctly updates the account balance.
func TestUpdateTransaction_AmountChange_UpdatesBalance(t *testing.T) {
	account := createTestAccountWithBalance(t, 1000)
	category := createTestCategory(t)

	// Create an expense of 300
	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Update Balance Test",
		Amount:          300,
		Currency:        "USD",
	}
	created, _, err := doCreateTransaction(createReq)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Balance should be 700 after expense of 300
	updatedAccount, _, _ := doGetAccount(account.Data.SubID)
	if updatedAccount.Data.Balance != 700 {
		t.Fatalf("Expected balance 700 after initial expense, got %d", updatedAccount.Data.Balance)
	}

	// Update the amount to 500
	newAmount := int64(500)
	updateReq := &models.UpdateTransactionRequest{
		Amount: &newAmount,
	}
	_, _, err = doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transaction: %v", err)
	}

	// Balance should now be 500 (1000 - 500)
	updatedAccount, _, _ = doGetAccount(account.Data.SubID)
	expectedBalance := int64(1000) - newAmount
	if updatedAccount.Data.Balance != expectedBalance {
		t.Errorf("Expected balance %d after updating expense to 500, got %d", expectedBalance, updatedAccount.Data.Balance)
	}
}

// TestUpdateTransaction_TypeChange_UpdatesBalanceCorrectly verifies that
// changing transaction type correctly adjusts the balance.
func TestUpdateTransaction_TypeChange_UpdatesBalanceCorrectly(t *testing.T) {
	account := createTestAccountWithBalance(t, 1000)
	category := createTestCategory(t)

	// Create an expense of 300 (balance will be 700)
	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Type Change Test",
		Amount:          300,
		Currency:        "USD",
	}
	created, _, err := doCreateTransaction(createReq)
	if err != nil {
		t.Fatalf("Failed to create expense transaction: %v", err)
	}

	// Verify balance after expense
	updatedAccount, _, _ := doGetAccount(account.Data.SubID)
	if updatedAccount.Data.Balance != 700 {
		t.Fatalf("Expected balance 700 after expense, got %d", updatedAccount.Data.Balance)
	}

	// Change from expense to income
	newType := models.TransactionTypeIncome
	updateReq := &models.UpdateTransactionRequest{
		TransactionType: &newType,
	}
	_, _, err = doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transaction type: %v", err)
	}

	// Balance should now be 1300:
	// Original 1000, reversed expense (+300) = 1300, applied income (+300) = 1300
	// But actually the logic is: reverse expense (add 300 back) = 1000, then apply income (add 300) = 1300
	updatedAccount, _, _ = doGetAccount(account.Data.SubID)
	expectedBalance := int64(1300)
	if updatedAccount.Data.Balance != expectedBalance {
		t.Errorf("Expected balance %d after changing expense to income, got %d", expectedBalance, updatedAccount.Data.Balance)
	}
}

// TestDeleteTransaction_RestoresBalance verifies that deleting a transaction
// restores the account balance to what it was before the transaction.
func TestDeleteTransaction_RestoresBalance(t *testing.T) {
	account := createTestAccountWithBalance(t, 1000)
	category := createTestCategory(t)

	// Create an expense of 400
	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Delete Balance Test",
		Amount:          400,
		Currency:        "USD",
	}
	created, _, err := doCreateTransaction(createReq)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Balance should be 600 after expense
	updatedAccount, _, _ := doGetAccount(account.Data.SubID)
	if updatedAccount.Data.Balance != 600 {
		t.Fatalf("Expected balance 600 after expense, got %d", updatedAccount.Data.Balance)
	}

	// Delete the transaction
	_, _, err = doDeleteTransaction(created.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to delete transaction: %v", err)
	}

	// Balance should be restored to 1000
	updatedAccount, _, _ = doGetAccount(account.Data.SubID)
	if updatedAccount.Data.Balance != 1000 {
		t.Errorf("Expected balance 1000 after delete, got %d", updatedAccount.Data.Balance)
	}
}

// TestDeleteTransfer_RestoresBothBalances verifies that deleting a transfer
// restores both source and destination account balances.
func TestDeleteTransfer_RestoresBothBalances(t *testing.T) {
	sourceAccount := createTestAccountWithBalance(t, 1000)
	destAccount := createTestAccountWithBalance(t, 500)
	category := createTestCategory(t)

	// Create a transfer of 300
	createReq := &models.CreateTransactionRequest{
		AccountID:         sourceAccount.Data.SubID,
		TransferAccountID: destAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   models.TransactionTypeTransfer,
		Title:             "Delete Transfer Test",
		Amount:            300,
		Currency:          "USD",
	}
	created, _, err := doCreateTransaction(createReq)
	if err != nil {
		t.Fatalf("Failed to create transfer transaction: %v", err)
	}

	// Verify balances after transfer
	sourceAfter, _, _ := doGetAccount(sourceAccount.Data.SubID)
	if sourceAfter.Data.Balance != 700 {
		t.Fatalf("Expected source balance 700 after transfer, got %d", sourceAfter.Data.Balance)
	}
	destAfter, _, _ := doGetAccount(destAccount.Data.SubID)
	if destAfter.Data.Balance != 800 {
		t.Fatalf("Expected dest balance 800 after transfer, got %d", destAfter.Data.Balance)
	}

	// Delete the transfer
	_, _, err = doDeleteTransaction(created.Data.SubID)
	if err != nil {
		t.Fatalf("Failed to delete transfer: %v", err)
	}

	// Both balances should be restored
	sourceAfter, _, _ = doGetAccount(sourceAccount.Data.SubID)
	if sourceAfter.Data.Balance != 1000 {
		t.Errorf("Expected source balance 1000 after delete, got %d", sourceAfter.Data.Balance)
	}
	destAfter, _, _ = doGetAccount(destAccount.Data.SubID)
	if destAfter.Data.Balance != 500 {
		t.Errorf("Expected dest balance 500 after delete, got %d", destAfter.Data.Balance)
	}
}

// TestUpdateExpenseToTransfer_UpdatesBalanceCorrectly verifies that updating
// an expense transaction to a transfer correctly reverses the original expense
// balance effect and applies new transfer balance effects to both accounts.
func TestUpdateExpenseToTransfer_UpdatesBalanceCorrectly(t *testing.T) {
	sourceAccount := createTestAccountWithBalance(t, 1000)
	destAccount := createTestAccountWithBalance(t, 500)
	category := createTestCategory(t)

	// Create an expense of 300 (source balance will be 700)
	createReq := &models.CreateTransactionRequest{
		AccountID:       sourceAccount.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Expense to Transfer Test",
		Amount:          300,
		Currency:        "USD",
	}
	created, _, err := doCreateTransaction(createReq)
	if err != nil {
		t.Fatalf("Failed to create expense transaction: %v", err)
	}

	// Verify balance after expense: 1000 - 300 = 700
	sourceAfter, _, _ := doGetAccount(sourceAccount.Data.SubID)
	if sourceAfter.Data.Balance != 700 {
		t.Fatalf("Expected source balance 700 after expense, got %d", sourceAfter.Data.Balance)
	}

	// Update expense to transfer with a new destination account
	transferAccountID := destAccount.Data.SubID
	updateReq := &models.UpdateTransactionRequest{
		TransactionType:   &[]models.TransactionType{models.TransactionTypeTransfer}[0],
		TransferAccountID: &transferAccountID,
	}
	_, _, err = doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update expense to transfer: %v", err)
	}

	// After update:
	// - Original expense is reversed: +300 to source (back to 1000)
	// - New transfer applied: -300 from source, +300 to dest
	// Expected: source = 700, dest = 800
	sourceAfter, _, _ = doGetAccount(sourceAccount.Data.SubID)
	expectedSourceBalance := int64(700)
	if sourceAfter.Data.Balance != expectedSourceBalance {
		t.Errorf("Expected source balance %d after expense->transfer, got %d", expectedSourceBalance, sourceAfter.Data.Balance)
	}

	destAfter, _, _ := doGetAccount(destAccount.Data.SubID)
	expectedDestBalance := int64(800)
	if destAfter.Data.Balance != expectedDestBalance {
		t.Errorf("Expected dest balance %d after expense->transfer, got %d", expectedDestBalance, destAfter.Data.Balance)
	}
}

// TestUpdateTransferToExpense_UpdatesBalanceCorrectly verifies that updating
// a transfer transaction to an expense correctly reverses the original transfer
// balance effects on both accounts and applies the new expense balance effect.
func TestUpdateTransferToExpense_UpdatesBalanceCorrectly(t *testing.T) {
	sourceAccount := createTestAccountWithBalance(t, 1000)
	destAccount := createTestAccountWithBalance(t, 500)
	category := createTestCategory(t)

	// Create a transfer of 300
	createReq := &models.CreateTransactionRequest{
		AccountID:         sourceAccount.Data.SubID,
		TransferAccountID: destAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   models.TransactionTypeTransfer,
		Title:             "Transfer to Expense Test",
		Amount:            300,
		Currency:          "USD",
	}
	created, _, err := doCreateTransaction(createReq)
	if err != nil {
		t.Fatalf("Failed to create transfer transaction: %v", err)
	}

	// Verify balances after transfer: source=700, dest=800
	sourceAfter, _, _ := doGetAccount(sourceAccount.Data.SubID)
	if sourceAfter.Data.Balance != 700 {
		t.Fatalf("Expected source balance 700 after transfer, got %d", sourceAfter.Data.Balance)
	}
	destAfter, _, _ := doGetAccount(destAccount.Data.SubID)
	if destAfter.Data.Balance != 800 {
		t.Fatalf("Expected dest balance 800 after transfer, got %d", destAfter.Data.Balance)
	}

	// Update transfer to expense (remove transfer_account_id)
	newType := models.TransactionTypeExpense
	updateReq := &models.UpdateTransactionRequest{
		TransactionType: &newType,
	}
	_, _, err = doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transfer to expense: %v", err)
	}

	// After update:
	// - Original transfer is reversed: +300 to source, -300 to dest (source=1000, dest=500)
	// - New expense applied: -300 from source (source=700)
	// Expected: source = 700, dest = 500
	sourceAfter, _, _ = doGetAccount(sourceAccount.Data.SubID)
	expectedSourceBalance := int64(700)
	if sourceAfter.Data.Balance != expectedSourceBalance {
		t.Errorf("Expected source balance %d after transfer->expense, got %d", expectedSourceBalance, sourceAfter.Data.Balance)
	}

	destAfter, _, _ = doGetAccount(destAccount.Data.SubID)
	expectedDestBalance := int64(500)
	if destAfter.Data.Balance != expectedDestBalance {
		t.Errorf("Expected dest balance %d after transfer->expense, got %d", expectedDestBalance, destAfter.Data.Balance)
	}
}

// TestCreateTransfer_InsufficientSourceBalance verifies that transfers are rejected
// when they would result in negative balance on the source account.
func TestCreateTransfer_InsufficientSourceBalance(t *testing.T) {
	sourceAccount := createTestAccountWithBalance(t, 100)
	destAccount := createTestAccountWithBalance(t, 500)
	category := createTestCategory(t)

	// Attempt to transfer 500 from account with only 100 balance
	transferAmount := int64(500)
	req := &models.CreateTransactionRequest{
		AccountID:         sourceAccount.Data.SubID,
		TransferAccountID: destAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   models.TransactionTypeTransfer,
		Title:             "Insufficient Balance Transfer",
		Amount:            transferAmount,
		Currency:          "USD",
	}
	_, status, _ := doCreateTransaction(req)

	// Transfer should be rejected when it would result in negative balance
	if status != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %d", status)
	}

	// Balances should be unchanged
	sourceAfter, _, _ := doGetAccount(sourceAccount.Data.SubID)
	if sourceAfter.Data.Balance != 100 {
		t.Errorf("Expected source balance 100 (unchanged), got %d", sourceAfter.Data.Balance)
	}
	destAfter, _, _ := doGetAccount(destAccount.Data.SubID)
	if destAfter.Data.Balance != 500 {
		t.Errorf("Expected dest balance 500 (unchanged), got %d", destAfter.Data.Balance)
	}
}

// createTestAccountWithBalance creates a test account with a specific initial balance.
func createTestAccountWithBalance(t *testing.T, initialBalance int64) *models.AccountResponse {
	accountReq := &models.CreateAccountRequest{
		Name:    "Balance Test Account",
		Type:    models.AccountTypeExpense,
		Balance: initialBalance,
	}
	account, _, err := doCreateAccount(accountReq)
	if err != nil {
		t.Fatalf("Failed to create test account: %v", err)
	}
	return account
}
