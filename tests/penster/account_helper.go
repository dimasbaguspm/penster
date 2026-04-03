package main

import "github.com/dimasbaguspm/penster/pkg/models"

// doCreateAccount POSTs a CreateAccountRequest and returns AccountResponse + status.
func doCreateAccount(req *models.CreateAccountRequest) (*models.AccountResponse, int, error) {
	result, status, err := doJSONRequest[models.AccountResponse]("POST", "/accounts", req)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doGetAccount GETs an account by ID and returns AccountResponse + status.
func doGetAccount(id string) (*models.AccountResponse, int, error) {
	result, status, err := doJSONRequest[models.AccountResponse]("GET", "/accounts/"+id, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doUpdateAccount PATCHes an account by ID and returns AccountResponse + status.
func doUpdateAccount(id string, req *models.UpdateAccountRequest) (*models.AccountResponse, int, error) {
	result, status, err := doJSONRequest[models.AccountResponse]("PATCH", "/accounts/"+id, req)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doListAccounts GETs the accounts list and returns AccountsResponse + status.
func doListAccounts() (*models.AccountsResponse, int, error) {
	result, status, err := doJSONRequest[models.AccountsResponse]("GET", "/accounts", nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}
