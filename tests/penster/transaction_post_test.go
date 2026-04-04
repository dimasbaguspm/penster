package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// createTestAccount creates a test account for transaction tests.
func createTestAccount(t *testing.T) *models.AccountResponse {
	accountReq := &models.CreateAccountRequest{
		Name:    "Test Account",
		Type:    models.AccountTypeExpense,
		Balance: 1000,
	}
	account, _, err := doCreateAccount(accountReq)
	if err != nil {
		t.Fatalf("Failed to create test account: %v", err)
	}
	return account
}

// createTestCategory creates a test category for transaction tests.
func createTestCategory(t *testing.T) *models.CategoryResponse {
	categoryReq := &models.CreateCategoryRequest{
		Name: "Test Category",
		Type: models.CategoryTypeExpense,
	}
	category, _, err := doCreateCategory(categoryReq)
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	return category
}

// createTransferTestAccount creates a second account for transfer transactions.
func createTransferTestAccount(t *testing.T) *models.AccountResponse {
	accountReq := &models.CreateAccountRequest{
		Name:    "Transfer Target Account",
		Type:    models.AccountTypeIncome,
		Balance: 0,
	}
	account, _, err := doCreateAccount(accountReq)
	if err != nil {
		t.Fatalf("Failed to create transfer target account: %v", err)
	}
	return account
}

func TestCreateTransaction_Success_Expense(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Test Expense",
		Amount:          500,
		Currency:        "USD",
	}
	result, status, err := doCreateTransaction(req)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Title != "Test Expense" {
		t.Errorf("Expected title 'Test Expense', got %s", result.Data.Title)
	}
	if result.Data.TransactionType != models.TransactionTypeExpense {
		t.Errorf("Expected type 'expense', got %s", result.Data.TransactionType)
	}
	if result.Data.Amount != 500 {
		t.Errorf("Expected amount 500, got %d", result.Data.Amount)
	}
	if result.Data.Currency != "USD" {
		t.Errorf("Expected currency 'USD', got %s", result.Data.Currency)
	}
}

func TestCreateTransaction_Success_Income(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeIncome,
		Title:           "Test Income",
		Amount:          1000,
		Currency:        "EUR",
	}
	result, status, err := doCreateTransaction(req)
	if err != nil {
		t.Fatalf("Failed to create income transaction: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.TransactionType != models.TransactionTypeIncome {
		t.Errorf("Expected type 'income', got %s", result.Data.TransactionType)
	}
	if result.Data.Amount != 1000 {
		t.Errorf("Expected amount 1000, got %d", result.Data.Amount)
	}
}

func TestCreateTransaction_Success_Transfer(t *testing.T) {
	account := createTestAccount(t)
	targetAccount := createTransferTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:         account.Data.SubID,
		TransferAccountID: targetAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   models.TransactionTypeTransfer,
		Title:             "Test Transfer",
		Amount:            250,
		Currency:          "USD",
	}
	result, status, err := doCreateTransaction(req)
	if err != nil {
		t.Fatalf("Failed to create transfer transaction: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.TransactionType != models.TransactionTypeTransfer {
		t.Errorf("Expected type 'transfer', got %s", result.Data.TransactionType)
	}
	if result.Data.TransferAccountID != targetAccount.Data.SubID {
		t.Errorf("Expected transfer_account_id '%s', got %s", targetAccount.Data.SubID, result.Data.TransferAccountID)
	}
}

func TestCreateTransaction_ValidationError_MissingAccountID(t *testing.T) {
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Missing Account",
		Amount:          100,
		Currency:        "USD",
	}
	_, status, _ := doCreateTransaction(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateTransaction_ValidationError_MissingCategoryID(t *testing.T) {
	account := createTestAccount(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Missing Category",
		Amount:          100,
		Currency:        "USD",
	}
	_, status, _ := doCreateTransaction(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateTransaction_ValidationError_MissingTitle(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Amount:          100,
		Currency:        "USD",
	}
	_, status, _ := doCreateTransaction(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateTransaction_ValidationError_MissingAmount(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Missing Amount",
		Currency:        "USD",
	}
	_, status, _ := doCreateTransaction(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateTransaction_ValidationError_MissingCurrency(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Missing Currency",
		Amount:          100,
	}
	_, status, _ := doCreateTransaction(req)
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateTransaction_ValidationError_InvalidTransactionType(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: "invalid_type",
		Title:           "Invalid Type",
		Amount:          100,
		Currency:        "USD",
	}
	status, err := doRequest("POST", "/transactions", req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
}

func TestCreateTransaction_WithNotes(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Transaction With Notes",
		Amount:          750,
		Currency:        "USD",
		Notes:           "This is a test note",
	}
	result, status, err := doCreateTransaction(req)
	if err != nil {
		t.Fatalf("Failed to create transaction with notes: %v", err)
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

func TestCreateTransaction_ZeroAmount(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Zero Amount Transaction",
		Amount:          0,
		Currency:        "USD",
	}
	result, status, err := doCreateTransaction(req)
	if err != nil {
		t.Fatalf("Failed to create zero amount transaction: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false, got true")
	}
}

func TestCreateTransaction_NegativeAmount(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	req := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Negative Amount Transaction",
		Amount:          -100,
		Currency:        "USD",
	}
	result, status, err := doCreateTransaction(req)
	if err != nil {
		t.Fatalf("Failed to create negative amount transaction: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for negative amount, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for negative amount, got true")
	}
}
