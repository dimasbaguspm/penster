package main

import (
	"fmt"
	"testing"

	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/response"
)

// doCreateDraft POSTs a CreateDraftRequest and returns DraftResponse + status.
func doCreateDraft(req *models.CreateDraftRequest) (*models.DraftResponse, int, error) {
	result, status, err := doJSONRequest[models.DraftResponse]("POST", "/drafts", req)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doGetDraft GETs a draft by ID and returns DraftResponse + status.
func doGetDraft(id string) (*models.DraftResponse, int, error) {
	result, status, err := doJSONRequest[models.DraftResponse]("GET", "/drafts/"+id, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doUpdateDraft PATCHes an UpdateDraftRequest and returns DraftResponse + status.
func doUpdateDraft(id string, req *models.UpdateDraftRequest) (*models.DraftResponse, int, error) {
	result, status, err := doJSONRequest[models.DraftResponse]("PATCH", "/drafts/"+id, req)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doListDrafts GETs the drafts list and returns DraftsResponse + status.
func doListDrafts() (*models.DraftsResponse, int, error) {
	result, status, err := doJSONRequest[models.DraftsResponse]("GET", "/drafts", nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doConfirmDraft POSTs to confirm a draft and returns TransactionResponse + status.
func doConfirmDraft(id string) (*models.TransactionResponse, int, error) {
	result, status, err := doJSONRequest[models.TransactionResponse]("POST", "/drafts/"+id+"/confirm", nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doRejectDraft POSTs to reject a draft and returns response.Response + status.
func doRejectDraft(id string) (*response.Response, int, error) {
	result, status, err := doJSONRequest[response.Response]("POST", "/drafts/"+id+"/reject", nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doDeleteDraft DELETEs a draft by ID and returns response.Response + status.
func doDeleteDraft(id string) (*response.Response, int, error) {
	result, status, err := doJSONRequest[response.Response]("DELETE", "/drafts/"+id, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// createTestDraftAccount creates a test account for draft tests.
func createTestDraftAccount(t *testing.T) *models.AccountResponse {
	accountReq := &models.CreateAccountRequest{
		Name:    "Draft Test Account",
		Type:    models.AccountTypeExpense,
		Balance: 1000,
	}
	account, _, err := doCreateAccount(accountReq)
	if err != nil {
		t.Fatalf("Failed to create test account: %v", err)
	}
	return account
}

// createTestDraftCategory creates a test category for draft tests.
func createTestDraftCategory(t *testing.T) *models.CategoryResponse {
	categoryReq := &models.CreateCategoryRequest{
		Name: "Draft Test Category",
		Type: models.CategoryTypeExpense,
	}
	category, _, err := doCreateCategory(categoryReq)
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	return category
}

// createTestDraftTransferAccount creates a second account for transfer drafts.
func createTestDraftTransferAccount(t *testing.T) *models.AccountResponse {
	accountReq := &models.CreateAccountRequest{
		Name:    "Draft Transfer Target Account",
		Type:    models.AccountTypeIncome,
		Balance: 0,
	}
	account, _, err := doCreateAccount(accountReq)
	if err != nil {
		t.Fatalf("Failed to create transfer target account: %v", err)
	}
	return account
}

// createTestDraftWithAccountAndCategory creates a basic draft with account and category.
func createTestDraftWithAccountAndCategory(t *testing.T) (*models.DraftResponse, *models.AccountResponse, *models.CategoryResponse) {
	account := createTestDraftAccount(t)
	category := createTestDraftCategory(t)

	// Use unique title to avoid test isolation issues
	uniqueTitle := fmt.Sprintf("Test Draft %s", account.Data.SubID)

	req := &models.CreateDraftRequest{
		AccountID:       account.Data.SubID,
		CategoryID:      category.Data.SubID,
		TransactionType: string(models.TransactionTypeExpense),
		Title:           uniqueTitle,
		Amount:          500,
		Currency:        "USD",
		TransactedAt:    "2024-01-15",
		Source:          string(models.DraftSourceManual),
	}
	draft, _, err := doCreateDraft(req)
	if err != nil {
		t.Fatalf("Failed to create test draft: %v", err)
	}
	return draft, account, category
}
