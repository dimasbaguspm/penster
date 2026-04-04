package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

func TestUpdateTransaction_PartialUpdate_Title(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Original Title",
		Amount:          500,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	newTitle := "Updated Title"
	updateReq := &models.UpdateTransactionRequest{
		Title: &newTitle,
	}

	result, status, err := doUpdateTransaction(id, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transaction title: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if !result.Success {
		t.Errorf("Expected success=true, got false with error: %s", result.Error)
	}
	if result.Data.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", result.Data.Title)
	}
	// Amount should remain unchanged
	if result.Data.Amount != 500 {
		t.Errorf("Expected amount unchanged 500, got %d", result.Data.Amount)
	}
}

func TestUpdateTransaction_PartialUpdate_Amount(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Amount Update Test",
		Amount:          100,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	newAmount := int64(999)
	updateReq := &models.UpdateTransactionRequest{
		Amount: &newAmount,
	}

	result, status, err := doUpdateTransaction(id, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transaction amount: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Amount != 999 {
		t.Errorf("Expected amount 999, got %d", result.Data.Amount)
	}
	// Title should remain unchanged
	if result.Data.Title != "Amount Update Test" {
		t.Errorf("Expected title unchanged 'Amount Update Test', got %s", result.Data.Title)
	}
}

func TestUpdateTransaction_PartialUpdate_Currency(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Currency Update Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	newCurrency := "EUR"
	updateReq := &models.UpdateTransactionRequest{
		Currency: &newCurrency,
	}

	result, status, err := doUpdateTransaction(id, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transaction currency: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Currency != "EUR" {
		t.Errorf("Expected currency 'EUR', got %s", result.Data.Currency)
	}
}

func TestUpdateTransaction_PartialUpdate_Notes(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Notes Update Test",
		Amount:          300,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	newNotes := "Updated notes"
	updateReq := &models.UpdateTransactionRequest{
		Notes: &newNotes,
	}

	result, status, err := doUpdateTransaction(id, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transaction notes: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Notes != "Updated notes" {
		t.Errorf("Expected notes 'Updated notes', got %s", result.Data.Notes)
	}
}

func TestUpdateTransaction_PartialUpdate_TransactionType(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Type Update Test",
		Amount:          400,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	newType := models.TransactionTypeIncome
	updateReq := &models.UpdateTransactionRequest{
		TransactionType: &newType,
	}

	result, status, err := doUpdateTransaction(id, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transaction type: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.TransactionType != models.TransactionTypeIncome {
		t.Errorf("Expected type 'income', got %s", result.Data.TransactionType)
	}
}

func TestUpdateTransaction_NotFound(t *testing.T) {
	nonExistentID := "00000000-0000-0000-0000-000000000000"
	newTitle := "Should Not Work"
	updateReq := &models.UpdateTransactionRequest{
		Title: &newTitle,
	}
	_, status, _ := doUpdateTransaction(nonExistentID, updateReq)
	if status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", status)
	}
}

func TestUpdateTransaction_InvalidUUID(t *testing.T) {
	newTitle := "Should Not Work"
	updateReq := &models.UpdateTransactionRequest{
		Title: &newTitle,
	}
	_, status, _ := doUpdateTransaction("not-a-valid-uuid", updateReq)
	if status != http.StatusBadRequest && status != http.StatusNotFound {
		t.Errorf("Expected status 400 or 404 for invalid UUID, got %d", status)
	}
}

func TestUpdateTransaction_AllFieldsNil(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Nil Update Test",
		Amount:          500,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	updateReq := &models.UpdateTransactionRequest{}

	result, status, err := doUpdateTransaction(id, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transaction with nil fields: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	// All fields should remain unchanged
	if result.Data.Title != "Nil Update Test" {
		t.Errorf("Expected title unchanged 'Nil Update Test', got %s", result.Data.Title)
	}
	if result.Data.Amount != 500 {
		t.Errorf("Expected amount unchanged 500, got %d", result.Data.Amount)
	}
	if result.Data.Currency != "USD" {
		t.Errorf("Expected currency unchanged 'USD', got %s", result.Data.Currency)
	}
}

func TestUpdateTransaction_MultipleFields(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Multiple Update Test",
		Amount:          100,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	newTitle := "New Multi Title"
	newAmount := int64(888)
	newCurrency := "GBP"
	updateReq := &models.UpdateTransactionRequest{
		Title:    &newTitle,
		Amount:   &newAmount,
		Currency: &newCurrency,
	}

	result, status, err := doUpdateTransaction(id, updateReq)
	if err != nil {
		t.Fatalf("Failed to update transaction multiple fields: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
	if result.Data.Title != "New Multi Title" {
		t.Errorf("Expected title 'New Multi Title', got %s", result.Data.Title)
	}
	if result.Data.Amount != 888 {
		t.Errorf("Expected amount 888, got %d", result.Data.Amount)
	}
	if result.Data.Currency != "GBP" {
		t.Errorf("Expected currency 'GBP', got %s", result.Data.Currency)
	}
}

func TestUpdateTransaction_Idempotent(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Idempotent Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)
	id := created.Data.SubID

	newTitle := "Idempotent Test"
	updateReq := &models.UpdateTransactionRequest{
		Title: &newTitle,
	}

	result1, status1, err1 := doUpdateTransaction(id, updateReq)
	if err1 != nil {
		t.Fatalf("First update failed: %v", err1)
	}
	if status1 != http.StatusOK {
		t.Errorf("Expected status 200 on first update, got %d", status1)
	}

	result2, status2, err2 := doUpdateTransaction(id, updateReq)
	if err2 != nil {
		t.Fatalf("Second update failed: %v", err2)
	}
	if status2 != http.StatusOK {
		t.Errorf("Expected status 200 on second update, got %d", status2)
	}

	if result1.Data.UpdatedAt.Unix() != result2.Data.UpdatedAt.Unix() {
		t.Logf("Note: UpdatedAt differs between calls (idempotent but timestamp may change)")
	}
}

func TestUpdateTransaction_InvalidAccountID(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Update Invalid Account Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	newAccountID := "00000000-0000-0000-0000-000000000000"
	updateReq := &models.UpdateTransactionRequest{
		AccountID: &newAccountID,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid account ID, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for invalid account ID, got true")
	}
}

func TestUpdateTransaction_InvalidCategoryID(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Update Invalid Category Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	newCategoryID := "00000000-0000-0000-0000-000000000000"
	updateReq := &models.UpdateTransactionRequest{
		CategoryID: &newCategoryID,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid category ID, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for invalid category ID, got true")
	}
}

func TestUpdateTransaction_InvalidTransferAccountID(t *testing.T) {
	account := createTestAccount(t)
	transferAccount := createTransferTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:         account.Data.SubID,
		TransferAccountID: transferAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   models.TransactionTypeTransfer,
		Title:            "Update Invalid Transfer Test",
		Amount:           200,
		Currency:         "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	newTransferAccountID := "00000000-0000-0000-0000-000000000000"
	updateReq := &models.UpdateTransactionRequest{
		TransferAccountID: &newTransferAccountID,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid transfer account ID, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for invalid transfer account ID, got true")
	}
}

func TestUpdateTransaction_MalformedAccountUUID(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Update Malformed UUID Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	malformedAccountID := "not-a-valid-uuid"
	updateReq := &models.UpdateTransactionRequest{
		AccountID: &malformedAccountID,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for malformed account UUID, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for malformed account UUID, got true")
	}
}

func TestUpdateTransaction_MalformedCategoryUUID(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Update Malformed Category UUID Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	malformedCategoryID := "not-a-valid-uuid"
	updateReq := &models.UpdateTransactionRequest{
		CategoryID: &malformedCategoryID,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for malformed category UUID, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for malformed category UUID, got true")
	}
}

func TestUpdateTransaction_TransferToSameAccount(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	// First create a transfer to a different account
	transferAccount := createTransferTestAccount(t)
	createReq := &models.CreateTransactionRequest{
		AccountID:         account.Data.SubID,
		TransferAccountID: transferAccount.Data.SubID,
		CategoryID:        category.Data.SubID,
		TransactionType:   models.TransactionTypeTransfer,
		Title:            "Original Transfer",
		Amount:           200,
		Currency:         "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	// Try to update transfer to point to the same account as source
	updateReq := &models.UpdateTransactionRequest{
		TransferAccountID: &account.Data.SubID,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for transfer to same account, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for transfer to same account, got true")
	}
}

func TestUpdateTransaction_EmptyAccountID(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Update Empty Account ID Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	emptyAccountID := ""
	updateReq := &models.UpdateTransactionRequest{
		AccountID: &emptyAccountID,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty account ID, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for empty account ID, got true")
	}
}

func TestUpdateTransaction_EmptyCategoryID(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Update Empty Category ID Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	emptyCategoryID := ""
	updateReq := &models.UpdateTransactionRequest{
		CategoryID: &emptyCategoryID,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty category ID, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for empty category ID, got true")
	}
}

func TestUpdateTransaction_NegativeAmount(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Update Negative Amount Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	negativeAmount := int64(-100)
	updateReq := &models.UpdateTransactionRequest{
		Amount: &negativeAmount,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for negative amount, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for negative amount, got true")
	}
}

func TestUpdateTransaction_ZeroAmount(t *testing.T) {
	account := createTestAccount(t)
	category := createTestCategory(t)

	createReq := &models.CreateTransactionRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: models.TransactionTypeExpense,
		Title:           "Update Zero Amount Test",
		Amount:          200,
		Currency:        "USD",
	}
	created, _, _ := doCreateTransaction(createReq)

	zeroAmount := int64(0)
	updateReq := &models.UpdateTransactionRequest{
		Amount: &zeroAmount,
	}

	result, status, err := doUpdateTransaction(created.Data.SubID, updateReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for zero amount, got %d", status)
	}
	if result.Success {
		t.Errorf("Expected success=false for zero amount, got true")
	}
}
