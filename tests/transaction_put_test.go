package main

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// Partial update table-driven tests
func TestUpdateTransaction_PartialUpdate_TableDriven(t *testing.T) {
	newTitle := "Updated Title"
	newAmount := int64(999)
	newCurrency := "EUR"
	newNotes := "Updated notes"
	newType := models.TransactionTypeIncome

	tests := []struct {
		name            string
		initialTitle    string
		initialAmount   int64
		initialCurrency string
		updateTitle     *string
		updateAmount    *int64
		updateCurrency  *string
		updateNotes     *string
		updateType      *models.TransactionType
		wantTitle       string
		wantAmount      int64
		wantCurrency    string
		wantNotes       string
	}{
		{
			name:            "update title only",
			initialTitle:    "Original Title",
			initialAmount:   500,
			initialCurrency: "USD",
			updateTitle:     &newTitle,
			wantTitle:       "Updated Title",
			wantAmount:      500,
			wantCurrency:    "USD",
		},
		{
			name:            "update amount only",
			initialTitle:    "Amount Update Test",
			initialAmount:   100,
			initialCurrency: "USD",
			updateAmount:    &newAmount,
			wantTitle:       "Amount Update Test",
			wantAmount:      999,
			wantCurrency:    "USD",
		},
		{
			name:            "update currency only",
			initialTitle:    "Currency Update Test",
			initialAmount:   200,
			initialCurrency: "USD",
			updateCurrency:  &newCurrency,
			wantTitle:       "Currency Update Test",
			wantAmount:      200,
			wantCurrency:    "EUR",
		},
		{
			name:            "update notes only",
			initialTitle:    "Notes Update Test",
			initialAmount:   300,
			initialCurrency: "USD",
			updateNotes:     &newNotes,
			wantTitle:       "Notes Update Test",
			wantAmount:      300,
			wantCurrency:    "USD",
			wantNotes:       "Updated notes",
		},
		{
			name:            "update transaction type only",
			initialTitle:    "Type Update Test",
			initialAmount:   400,
			initialCurrency: "USD",
			updateType:      &newType,
			wantTitle:       "Type Update Test",
			wantAmount:      400,
			wantCurrency:    "USD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := createTestAccount(t)
			category := createTestCategory(t)

			createReq := &models.CreateTransactionRequest{
				AccountID:       account.Data.SubID,
				CategoryID:      category.Data.SubID,
				TransactionType: models.TransactionTypeExpense,
				Title:           tt.initialTitle,
				Amount:          tt.initialAmount,
				Currency:        tt.initialCurrency,
			}
			created, _, _ := doCreateTransaction(createReq)
			id := created.Data.SubID

			updateReq := &models.UpdateTransactionRequest{
				Title:           tt.updateTitle,
				Amount:          tt.updateAmount,
				Currency:        tt.updateCurrency,
				Notes:           tt.updateNotes,
				TransactionType: tt.updateType,
			}

			result, status, err := doUpdateTransaction(id, updateReq)
			if err != nil {
				t.Fatalf("Failed to update transaction: %v", err)
			}
			if status != http.StatusOK {
				t.Errorf("Expected status 200, got %d", status)
			}
			if tt.updateTitle != nil && result.Data.Title != tt.wantTitle {
				t.Errorf("Expected title %q, got %q", tt.wantTitle, result.Data.Title)
			}
			if tt.updateAmount != nil && result.Data.Amount != tt.wantAmount {
				t.Errorf("Expected amount %d, got %d", tt.wantAmount, result.Data.Amount)
			}
			if tt.updateCurrency != nil && result.Data.Currency != tt.wantCurrency {
				t.Errorf("Expected currency %q, got %q", tt.wantCurrency, result.Data.Currency)
			}
			if tt.updateNotes != nil && result.Data.Notes != tt.wantNotes {
				t.Errorf("Expected notes %q, got %q", tt.wantNotes, result.Data.Notes)
			}
			if tt.updateType != nil && result.Data.TransactionType != *tt.updateType {
				t.Errorf("Expected type %q, got %q", *tt.updateType, result.Data.TransactionType)
			}
		})
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

// Validation table-driven tests
func TestUpdateTransaction_Validation_TableDriven(t *testing.T) {
	nonExistentID := "00000000-0000-0000-0000-000000000000"
	malformedID := "not-a-valid-uuid"
	negativeAmount := int64(-100)
	zeroAmount := int64(0)

	tests := []struct {
		name      string
		setup     func(t *testing.T) (id string, updateReq *models.UpdateTransactionRequest)
		wantStatus int
	}{
		{
			name: "invalid account ID",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:       account.Data.SubID,
					CategoryID:      category.Data.SubID,
					TransactionType: models.TransactionTypeExpense,
					Title:           "Test",
					Amount:          200,
					Currency:        "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				updateReq := &models.UpdateTransactionRequest{
					AccountID: &nonExistentID,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid category ID",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:       account.Data.SubID,
					CategoryID:      category.Data.SubID,
					TransactionType: models.TransactionTypeExpense,
					Title:           "Test",
					Amount:          200,
					Currency:        "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				updateReq := &models.UpdateTransactionRequest{
					CategoryID: &nonExistentID,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid transfer account ID",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				transferAccount := createTransferTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:         account.Data.SubID,
					TransferAccountID: transferAccount.Data.SubID,
					CategoryID:        category.Data.SubID,
					TransactionType:   models.TransactionTypeTransfer,
					Title:             "Test",
					Amount:            200,
					Currency:          "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				updateReq := &models.UpdateTransactionRequest{
					TransferAccountID: &nonExistentID,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "malformed account UUID",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:       account.Data.SubID,
					CategoryID:      category.Data.SubID,
					TransactionType: models.TransactionTypeExpense,
					Title:           "Test",
					Amount:          200,
					Currency:        "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				updateReq := &models.UpdateTransactionRequest{
					AccountID: &malformedID,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "malformed category UUID",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:       account.Data.SubID,
					CategoryID:      category.Data.SubID,
					TransactionType: models.TransactionTypeExpense,
					Title:           "Test",
					Amount:          200,
					Currency:        "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				updateReq := &models.UpdateTransactionRequest{
					CategoryID: &malformedID,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "transfer to same account",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				transferAccount := createTransferTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:         account.Data.SubID,
					TransferAccountID: transferAccount.Data.SubID,
					CategoryID:        category.Data.SubID,
					TransactionType:   models.TransactionTypeTransfer,
					Title:             "Test",
					Amount:            200,
					Currency:          "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				updateReq := &models.UpdateTransactionRequest{
					TransferAccountID: &account.Data.SubID,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty account ID",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:       account.Data.SubID,
					CategoryID:      category.Data.SubID,
					TransactionType: models.TransactionTypeExpense,
					Title:           "Test",
					Amount:          200,
					Currency:        "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				emptyID := ""
				updateReq := &models.UpdateTransactionRequest{
					AccountID: &emptyID,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty category ID",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:       account.Data.SubID,
					CategoryID:      category.Data.SubID,
					TransactionType: models.TransactionTypeExpense,
					Title:           "Test",
					Amount:          200,
					Currency:        "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				emptyID := ""
				updateReq := &models.UpdateTransactionRequest{
					CategoryID: &emptyID,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative amount",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:       account.Data.SubID,
					CategoryID:      category.Data.SubID,
					TransactionType: models.TransactionTypeExpense,
					Title:           "Test",
					Amount:          200,
					Currency:        "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				updateReq := &models.UpdateTransactionRequest{
					Amount: &negativeAmount,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "zero amount",
			setup: func(t *testing.T) (string, *models.UpdateTransactionRequest) {
				account := createTestAccount(t)
				category := createTestCategory(t)
				createReq := &models.CreateTransactionRequest{
					AccountID:       account.Data.SubID,
					CategoryID:      category.Data.SubID,
					TransactionType: models.TransactionTypeExpense,
					Title:           "Test",
					Amount:          200,
					Currency:        "USD",
				}
				created, _, _ := doCreateTransaction(createReq)
				updateReq := &models.UpdateTransactionRequest{
					Amount: &zeroAmount,
				}
				return created.Data.SubID, updateReq
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, updateReq := tt.setup(t)

			_, status, err := doUpdateTransaction(id, updateReq)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			if status != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, status)
			}
		})
	}
}
