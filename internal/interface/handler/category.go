package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/interface/dto"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/response"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// List handles GET /categories
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := dto.ParseCategoryListParams(r)
	categories, total, err := h.svc.List(r.Context(), params)
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
func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	category, err := h.svc.GetByID(r.Context(), id)
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
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	category, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := response.NewResponse(*category)
	h.writeJSON(w, http.StatusCreated, resp)
}

// Update handles PUT /categories/:id
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	category, err := h.svc.Update(r.Context(), id, &req)
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
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	category, err := h.svc.Delete(r.Context(), id)
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
