package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/interface/dto"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
	"github.com/dimasbaguspm/penster/pkg/response"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// List handles GET /categories
// @Summary List all categories
// @Description Get a paginated list of categories with optional filtering
// @Tags categories
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param sort_by query string false "Sort by field"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.Response
// @Router /categories [get]
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Category", "List")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := dto.ParseCategoryListParams(r)
	categories, total, err := h.svc.List(ctx, params)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	categoryList := make([]models.Category, 0, len(categories))
	for _, cat := range categories {
		categoryList = append(categoryList, *cat)
	}

	resp := response.NewPaginatedResponse(categoryList, params.PageNumber, params.PageSize, total)
	h.writeJSON(w, http.StatusOK, resp)
}

// Get handles GET /categories/:id
// @Summary Get category by ID
// @Description Get a single category by its UUID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{id} [get]
func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Category", "Get")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	category, err := h.svc.GetByID(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if category == nil {
		h.writeError(w, http.StatusNotFound, "category not found")
		return
	}

	resp := response.NewResponse(*category)
	h.writeJSON(w, http.StatusOK, resp)
}

// Create handles POST /categories
// @Summary Create a new category
// @Description Create a new category with name and type
// @Tags categories
// @Accept json
// @Produce json
// @Param request body models.CreateCategoryRequest true "Category creation request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Category", "Create")
	defer span.End()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateCreateCategoryRequest(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.svc.Create(ctx, &req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := response.NewResponse(*category)
	h.writeJSON(w, http.StatusCreated, resp)
}

// Update handles PUT /categories/:id
// @Summary Update a category
// @Description Update an existing category by its UUID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category UUID"
// @Param request body models.UpdateCategoryRequest true "Category update request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Category", "Update")
	defer span.End()

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	var req models.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateUpdateCategoryRequest(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.svc.Update(ctx, id, &req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if category == nil {
		h.writeError(w, http.StatusNotFound, "category not found")
		return
	}

	resp := response.NewResponse(*category)
	h.writeJSON(w, http.StatusOK, resp)
}

// Delete handles DELETE /categories/:id
// @Summary Delete a category
// @Description Soft delete a category by its UUID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Category", "Delete")
	defer span.End()

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	category, err := h.svc.Delete(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if category == nil {
		h.writeError(w, http.StatusNotFound, "category not found")
		return
	}

	resp := response.NewResponse(*category)
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *CategoryHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *CategoryHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.NewErrorResponse(message))
}
