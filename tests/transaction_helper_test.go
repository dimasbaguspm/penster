package main

import (
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
)

// doCreateTransaction POSTs a CreateTransactionRequest and returns TransactionResponse + status.
func doCreateTransaction(req *models.CreateTransactionRequest) (*models.TransactionResponse, int, error) {
	result, status, err := doJSONRequest[models.TransactionResponse]("POST", "/transactions", req)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doGetTransaction GETs a transaction by ID and returns TransactionResponse + status.
func doGetTransaction(id string) (*models.TransactionResponse, int, error) {
	result, status, err := doJSONRequest[models.TransactionResponse]("GET", "/transactions/"+id, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doUpdateTransaction PUTs an UpdateTransactionRequest and returns TransactionResponse + status.
func doUpdateTransaction(id string, req *models.UpdateTransactionRequest) (*models.TransactionResponse, int, error) {
	result, status, err := doJSONRequest[models.TransactionResponse]("PUT", "/transactions/"+id, req)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doDeleteTransaction DELETEs a transaction by ID and returns ErrorResponse + status.
func doDeleteTransaction(id string) (*models.ErrorResponse, int, error) {
	result, status, err := doJSONRequest[models.ErrorResponse]("DELETE", "/transactions/"+id, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doListTransactions GETs the transactions list and returns TransactionPagedResponse + status.
func doListTransactions() (*models.TransactionPagedResponse, int, error) {
	result, status, err := doJSONRequest[models.TransactionPagedResponse]("GET", "/transactions", nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

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
