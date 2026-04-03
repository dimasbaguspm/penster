package main

import "github.com/dimasbaguspm/penster/pkg/models"

// doCreateCategory POSTs a CreateCategoryRequest and returns CategoryResponse + status.
func doCreateCategory(req *models.CreateCategoryRequest) (*models.CategoryResponse, int, error) {
	result, status, err := doJSONRequest[models.CategoryResponse]("POST", "/categories", req)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doGetCategory GETs a category by ID and returns CategoryResponse + status.
func doGetCategory(id string) (*models.CategoryResponse, int, error) {
	result, status, err := doJSONRequest[models.CategoryResponse]("GET", "/categories/"+id, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doUpdateCategory PUTs an UpdateCategoryRequest and returns CategoryResponse + status.
func doUpdateCategory(id string, req *models.UpdateCategoryRequest) (*models.CategoryResponse, int, error) {
	result, status, err := doJSONRequest[models.CategoryResponse]("PUT", "/categories/"+id, req)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doDeleteCategory DELETEs a category by ID and returns CategoryResponse + status.
func doDeleteCategory(id string) (*models.CategoryResponse, int, error) {
	result, status, err := doJSONRequest[models.CategoryResponse]("DELETE", "/categories/"+id, nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}

// doListCategories GETs the categories list and returns CategoriesResponse + status.
func doListCategories() (*models.CategoriesResponse, int, error) {
	result, status, err := doJSONRequest[models.CategoriesResponse]("GET", "/categories", nil)
	if result == nil {
		return nil, status, err
	}
	return result, status, err
}
