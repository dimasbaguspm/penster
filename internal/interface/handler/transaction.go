package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/domain/entities"
	"github.com/dimasbaguspm/penster/internal/interface/dto"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
	"github.com/dimasbaguspm/penster/pkg/response"
)

type TransactionHandler struct {
	svc *service.TransactionService
}

func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// List handles GET /transactions
// @Summary List all transactions
// @Description Get a paginated list of transactions with optional filtering
// @Tags transactions
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param account_id query string false "Filter by account ID"
// @Param category_id query string false "Filter by category ID"
// @Param transaction_type query string false "Filter by transaction type (expense, income, transfer)"
// @Param sort_by query string false "Sort by field (title, transacted_at, created_at, amount)"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.Response
// @Router /transactions [get]
func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Transaction", "List")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := dto.ParseTransactionListParams(r)
	transactions, total, err := h.svc.List(ctx, params)

	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	transactionList := make([]models.Transaction, 0, len(transactions))
	for _, tx := range transactions {
		transactionList = append(transactionList, *tx)
	}

	resp := response.NewPaginatedResponse(transactionList, params.PageNumber, params.PageSize, total)
	h.writeJSON(w, http.StatusOK, resp)
}

// Get handles GET /transactions/:id
// @Summary Get transaction by ID
// @Description Get a single transaction by its UUID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /transactions/{id} [get]
func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Transaction", "Get")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	tx, err := h.svc.GetByID(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tx == nil {
		h.writeError(w, http.StatusNotFound, "transaction not found")
		return
	}

	resp := response.NewResponse(*tx)
	h.writeJSON(w, http.StatusOK, resp)
}

// Create handles POST /transactions
// @Summary Create a new transaction
// @Description Create a new transaction with amount, type, and account/category references
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body models.CreateTransactionRequest true "Transaction creation request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /transactions [post]
func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Transaction", "Create")
	defer span.End()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateCreateTransactionRequest(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	tx, err := h.svc.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, entities.ErrAccountNotFound) ||
			errors.Is(err, entities.ErrCategoryNotFound) ||
			errors.Is(err, entities.ErrTransferAccountNotFound) ||
			errors.Is(err, entities.ErrInsufficientBalance) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := response.NewResponse(*tx)
	h.writeJSON(w, http.StatusCreated, resp)
}

// Update handles PUT /transactions/:id
// @Summary Update a transaction
// @Description Update an existing transaction by its UUID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction UUID"
// @Param request body models.UpdateTransactionRequest true "Transaction update request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /transactions/{id} [put]
func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Transaction", "Update")
	defer span.End()

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	var req models.UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateUpdateTransactionRequest(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	tx, err := h.svc.Update(ctx, id, &req)
	if err != nil {
		if errors.Is(err, entities.ErrAccountNotFound) ||
			errors.Is(err, entities.ErrCategoryNotFound) ||
			errors.Is(err, entities.ErrTransferAccountNotFound) ||
			errors.Is(err, entities.ErrTransferToSameAccount) ||
			errors.Is(err, entities.ErrInsufficientBalance) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tx == nil {
		h.writeError(w, http.StatusNotFound, "transaction not found")
		return
	}

	resp := response.NewResponse(*tx)
	h.writeJSON(w, http.StatusOK, resp)
}

// Delete handles DELETE /transactions/:id
// @Summary Delete a transaction
// @Description Soft delete a transaction by its UUID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Transaction", "Delete")
	defer span.End()

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	tx, err := h.svc.Delete(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tx == nil {
		h.writeError(w, http.StatusNotFound, "transaction not found")
		return
	}

	resp := response.NewResponse(*tx)
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *TransactionHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *TransactionHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.NewErrorResponse(message))
}
