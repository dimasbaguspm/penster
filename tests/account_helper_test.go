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

// doUpdateAccount PUTs an UpdateAccountRequest and returns AccountResponse + status.
func doUpdateAccount(id string, req *models.UpdateAccountRequest) (*models.AccountResponse, int, error) {
	result, status, err := doJSONRequest[models.AccountResponse]("PUT", "/accounts/"+id, req)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doDeleteAccount DELETEs an account by ID and returns AccountResponse + status.
func doDeleteAccount(id string) (*models.AccountResponse, int, error) {
	result, status, err := doJSONRequest[models.AccountResponse]("DELETE", "/accounts/"+id, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doListAccounts GETs the accounts list and returns AccountPagedResponse + status.
func doListAccounts() (*models.AccountPagedResponse, int, error) {
	result, status, err := doJSONRequest[models.AccountPagedResponse]("GET", "/accounts", nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}
