package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

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

func TestCreateTransaction_ValidationError_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		setupReq  func(t *testing.T, req *models.CreateTransactionRequest)
		wantStatus int
	}{
		{
			name: "missing_account_id",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				category := createTestCategory(t)
				req.CategoryID = category.Data.SubID
				req.AccountID = ""
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing_category_id",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				req.AccountID = account.Data.SubID
				req.CategoryID = ""
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing_title",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				req.AccountID = account.Data.SubID
				req.CategoryID = category.Data.SubID
				req.Title = ""
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing_amount",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				req.AccountID = account.Data.SubID
				req.CategoryID = category.Data.SubID
				req.Amount = 0
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing_currency",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				req.AccountID = account.Data.SubID
				req.CategoryID = category.Data.SubID
				req.Currency = ""
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_account_id",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				category := createTestCategory(t)
				req.AccountID = "00000000-0000-0000-0000-000000000000"
				req.CategoryID = category.Data.SubID
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_category_id",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				req.AccountID = account.Data.SubID
				req.CategoryID = "00000000-0000-0000-0000-000000000000"
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_transfer_account_id",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				req.AccountID = account.Data.SubID
				req.TransferAccountID = "00000000-0000-0000-0000-000000000000"
				req.CategoryID = category.Data.SubID
				req.TransactionType = models.TransactionTypeTransfer
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "malformed_account_uuid",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				category := createTestCategory(t)
				req.AccountID = "not-a-valid-uuid"
				req.CategoryID = category.Data.SubID
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "malformed_category_uuid",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				req.AccountID = account.Data.SubID
				req.CategoryID = "not-a-valid-uuid"
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "transfer_to_same_account",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				req.AccountID = account.Data.SubID
				req.TransferAccountID = account.Data.SubID
				req.CategoryID = category.Data.SubID
				req.TransactionType = models.TransactionTypeTransfer
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty_account_id",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				category := createTestCategory(t)
				req.AccountID = ""
				req.CategoryID = category.Data.SubID
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty_category_id",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				req.AccountID = account.Data.SubID
				req.CategoryID = ""
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "zero_amount",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				req.AccountID = account.Data.SubID
				req.CategoryID = category.Data.SubID
				req.Amount = 0
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_amount",
			setupReq: func(t *testing.T, req *models.CreateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				req.AccountID = account.Data.SubID
				req.CategoryID = category.Data.SubID
				req.Amount = -100
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &models.CreateTransactionRequest{
				TransactionType: models.TransactionTypeExpense,
				Title:           "Test Transaction",
				Amount:          100,
				Currency:        "USD",
			}
			tt.setupReq(t, req)

			result, status, err := doCreateTransaction(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			if status != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, status)
			}
			if result.Success {
				t.Errorf("Expected success=false, got true")
			}
		})
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
