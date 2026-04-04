package main

import "github.com/dimasbaguspm/penster/pkg/models"

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

// doDeleteTransaction DELETEs a transaction by ID and returns TransactionResponse + status.
func doDeleteTransaction(id string) (*models.TransactionResponse, int, error) {
	result, status, err := doJSONRequest[models.TransactionResponse]("DELETE", "/transactions/"+id, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doListTransactions GETs the transactions list and returns TransactionsResponse + status.
func doListTransactions() (*models.TransactionsResponse, int, error) {
	result, status, err := doJSONRequest[models.TransactionsResponse]("GET", "/transactions", nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}
