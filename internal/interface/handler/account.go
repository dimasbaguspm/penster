package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/interface/dto"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/response"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

// List handles GET /accounts
// @Summary List all accounts
// @Description Get a paginated list of accounts with optional filtering
// @Tags accounts
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param sort_by query string false "Sort by field"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.Response
// @Router /accounts [get]
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := dto.ParseAccountListParams(r)
	accounts, total, err := h.svc.List(r.Context(), params)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	accountList := make([]models.Account, 0, len(accounts))
	for _, acc := range accounts {
		accountList = append(accountList, *acc)
	}

	resp := response.NewPaginatedResponse(accountList, params.PageNumber, params.PageSize, total)
	h.writeJSON(w, http.StatusOK, resp)
}

// Get handles GET /accounts/:id
// @Summary Get account by ID
// @Description Get a single account by its UUID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /accounts/{id} [get]
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if account == nil {
		h.writeError(w, http.StatusNotFound, "account not found")
		return
	}

	resp := response.NewResponse(*account)
	h.writeJSON(w, http.StatusOK, resp)
}

// Create handles POST /accounts
// @Summary Create a new account
// @Description Create a new account with name, type and optional balance
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body models.CreateAccountRequest true "Account creation request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /accounts [post]
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateCreateAccountRequest(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := response.NewResponse(*account)
	h.writeJSON(w, http.StatusCreated, resp)
}

// Update handles PUT /accounts/:id
// @Summary Update an account
// @Description Update an existing account by its UUID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account UUID"
// @Param request body models.UpdateAccountRequest true "Account update request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /accounts/{id} [put]
func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var req models.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateUpdateAccountRequest(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if account == nil {
		h.writeError(w, http.StatusNotFound, "account not found")
		return
	}

	resp := response.NewResponse(*account)
	h.writeJSON(w, http.StatusOK, resp)
}

// Delete handles DELETE /accounts/:id
// @Summary Delete an account
// @Description Soft delete an account by its UUID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /accounts/{id} [delete]
func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := h.svc.Delete(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if account == nil {
		h.writeError(w, http.StatusNotFound, "account not found")
		return
	}

	resp := response.NewResponse(*account)
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *AccountHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AccountHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.NewErrorResponse(message))
}
