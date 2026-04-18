package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/interface/dto"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
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
// @Success 200 {object} models.AccountPagedResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts [get]
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "account")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Account", "List")
	defer span.End()

	if r.Method != http.MethodGet {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Info("listing accounts")
	params := dto.ParseAccountListParams(r)
	accounts, total, err := h.svc.List(ctx, params)
	if err != nil {
		log.Error("failed to list accounts", "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	accountList := make([]models.Account, 0, len(accounts))
	for _, acc := range accounts {
		accountList = append(accountList, *acc)
	}

	log.Info("accounts listed", "count", len(accountList), "total", total)
	resp := models.NewAccountPagedResponse(accountList, params.PageSize, params.PageNumber, total)
	h.writeJSON(w, http.StatusOK, resp)
}

// Get handles GET /accounts/:id
// @Summary Get account by ID
// @Description Get a single account by its UUID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account UUID"
// @Success 200 {object} models.AccountResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts/{id} [get]
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "account")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Account", "Get")
	defer span.End()

	if r.Method != http.MethodGet {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		log.Warn("invalid account id")
		h.writeError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	log.Info("getting account", "id", id)
	account, err := h.svc.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to get account", "id", id, "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if account == nil {
		log.Info("account not found", "id", id)
		h.writeError(w, http.StatusNotFound, "account not found")
		return
	}

	log.Info("account retrieved", "id", id)
	resp := models.AccountResponse{Data: *account}
	h.writeJSON(w, http.StatusOK, resp)
}

// Create handles POST /accounts
// @Summary Create a new account
// @Description Create a new account with name, type and optional balance
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body models.CreateAccountRequest true "Account creation request"
// @Success 201 {object} models.AccountResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts [post]
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "account")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Account", "Create")
	defer span.End()

	if r.Method != http.MethodPost {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid request body", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateCreateAccountRequest(ctx, &req); err != nil {
		log.Warn("validation failed", "error", err)
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Info("creating account", "name", req.Name, "type", req.Type)
	account, err := h.svc.Create(ctx, &req)
	if err != nil {
		log.Error("failed to create account", "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("account created", "id", account.ID)
	resp := models.AccountResponse{Data: *account}
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
// @Success 200 {object} models.AccountResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts/{id} [put]
func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "account")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Account", "Update")
	defer span.End()

	if r.Method != http.MethodPut {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		log.Warn("invalid account id")
		h.writeError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var req models.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid request body", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateUpdateAccountRequest(ctx, &req); err != nil {
		log.Warn("validation failed", "error", err)
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Info("updating account", "id", id)
	account, err := h.svc.Update(ctx, id, &req)
	if err != nil {
		log.Error("failed to update account", "id", id, "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if account == nil {
		log.Info("account not found", "id", id)
		h.writeError(w, http.StatusNotFound, "account not found")
		return
	}

	log.Info("account updated", "id", id)
	resp := models.AccountResponse{Data: *account}
	h.writeJSON(w, http.StatusOK, resp)
}

// Delete handles DELETE /accounts/:id
// @Summary Delete an account
// @Description Soft delete an account by its UUID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account UUID"
// @Success 200 {object} models.AccountResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts/{id} [delete]
func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "account")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Account", "Delete")
	defer span.End()

	if r.Method != http.MethodDelete {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		log.Warn("invalid account id")
		h.writeError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	log.Info("deleting account", "id", id)
	account, err := h.svc.Delete(ctx, id)
	if err != nil {
		log.Error("failed to delete account", "id", id, "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if account == nil {
		log.Info("account not found", "id", id)
		h.writeError(w, http.StatusNotFound, "account not found")
		return
	}

	log.Info("account deleted", "id", id)
	resp := models.AccountResponse{Data: *account}
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
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}
