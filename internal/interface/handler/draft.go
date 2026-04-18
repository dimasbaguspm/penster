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
	"github.com/google/uuid"
)

type DraftHandler struct {
	svc *service.DraftService
}

func NewDraftHandler(svc *service.DraftService) *DraftHandler {
	return &DraftHandler{svc: svc}
}

// List handles GET /drafts
// @Summary List all drafts
// @Description Get a paginated list of drafts with optional filtering by source and status
// @Tags drafts
// @Accept json
// @Produce json
// @Param source query string false "Filter by source (manual, ingestion)"
// @Param status query string false "Filter by status (pending, confirmed, rejected)"
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.DraftPagedResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /drafts [get]
func (h *DraftHandler) List(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "draft")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Draft", "List")
	defer span.End()

	if r.Method != http.MethodGet {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Info("listing drafts")
	params := dto.ParseDraftListParams(r)

	drafts, total, err := h.svc.List(ctx, params)
	if err != nil {
		log.Error("failed to list drafts", "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	draftList := make([]models.Draft, 0, len(drafts))
	for _, d := range drafts {
		draftList = append(draftList, *d)
	}

	log.Info("drafts listed", "count", len(draftList), "total", total)
	resp := models.NewDraftPagedResponse(draftList, params.PageSize, 1, total)
	h.writeJSON(w, http.StatusOK, resp)
}

// Get handles GET /drafts/:id
// @Summary Get draft by ID
// @Description Get a single draft by its UUID
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft UUID"
// @Success 200 {object} models.DraftResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /drafts/{id} [get]
func (h *DraftHandler) Get(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "draft")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Draft", "Get")
	defer span.End()

	if r.Method != http.MethodGet {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		log.Warn("invalid draft id")
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid draft id format", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	log.Info("getting draft", "id", id)
	draft, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			log.Info("draft not found", "id", id)
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Error("failed to get draft", "id", id, "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if draft == nil {
		log.Info("draft not found", "id", id)
		h.writeError(w, http.StatusNotFound, "draft not found")
		return
	}

	log.Info("draft retrieved", "id", id)
	resp := models.DraftResponse{Data: *draft}
	h.writeJSON(w, http.StatusOK, resp)
}

// Create handles POST /drafts
// @Summary Create a new draft
// @Description Create a new draft with transaction details
// @Tags drafts
// @Accept json
// @Produce json
// @Param request body models.CreateDraftRequest true "Draft creation request"
// @Success 201 {object} models.DraftResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /drafts [post]
func (h *DraftHandler) Create(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "draft")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Draft", "Create")
	defer span.End()

	if r.Method != http.MethodPost {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid request body", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateCreateDraftRequest(ctx, &req); err != nil {
		log.Warn("validation failed", "error", err)
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Info("creating draft")
	draft, err := h.svc.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, entities.ErrAccountNotFound) ||
			errors.Is(err, entities.ErrCategoryNotFound) ||
			errors.Is(err, entities.ErrTransferAccountNotFound) {
			log.Warn("bad request creating draft", "error", err)
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Error("failed to create draft", "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("draft created", "id", draft.ID)
	resp := models.DraftResponse{Data: *draft}
	h.writeJSON(w, http.StatusCreated, resp)
}

// Update handles PATCH /drafts/:id
// @Summary Update a draft
// @Description Update an existing draft by its UUID
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft UUID"
// @Param request body models.UpdateDraftRequest true "Draft update request"
// @Success 200 {object} models.DraftResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /drafts/{id} [patch]
func (h *DraftHandler) Update(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "draft")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Draft", "Update")
	defer span.End()

	if r.Method != http.MethodPatch {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		log.Warn("invalid draft id")
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid draft id format", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	var req models.UpdateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid request body", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateUpdateDraftRequest(ctx, &req); err != nil {
		log.Warn("validation failed", "error", err)
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Info("updating draft", "id", id)
	draft, err := h.svc.Update(ctx, id, &req)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			log.Info("draft not found", "id", id)
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, entities.ErrAccountNotFound) ||
			errors.Is(err, entities.ErrCategoryNotFound) ||
			errors.Is(err, entities.ErrTransferAccountNotFound) {
			log.Warn("bad request updating draft", "error", err)
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Error("failed to update draft", "id", id, "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Validate same-account transfer after getting current draft
	if req.TransferAccountID != nil && req.AccountID != nil &&
		*req.TransferAccountID == *req.AccountID {
		log.Warn("transfer account same as source account", "id", id)
		h.writeError(w, http.StatusBadRequest, "transfer account cannot be the same as source account")
		return
	}

	log.Info("draft updated", "id", id)
	resp := models.DraftResponse{Data: *draft}
	h.writeJSON(w, http.StatusOK, resp)
}

// Confirm handles POST /drafts/:id/confirm
// @Summary Confirm a draft
// @Description Promote a draft to a committed transaction
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft UUID"
// @Success 200 {object} models.TransactionResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /drafts/{id}/confirm [post]
func (h *DraftHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "draft")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Draft", "Confirm")
	defer span.End()

	if r.Method != http.MethodPost {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		log.Warn("invalid draft id")
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid draft id format", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	log.Info("confirming draft", "id", id)
	tx, err := h.svc.Confirm(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			log.Info("draft not found", "id", id)
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, entities.ErrDraftNotPending) {
			log.Warn("draft not pending", "id", id, "error", err)
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, entities.ErrInsufficientBalance) {
			log.Warn("insufficient balance", "id", id, "error", err)
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Error("failed to confirm draft", "id", id, "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("draft confirmed", "id", id, "transaction_id", tx.ID)
	resp := models.TransactionResponse{Data: *tx}
	h.writeJSON(w, http.StatusOK, resp)
}

// Reject handles POST /drafts/:id/reject
// @Summary Reject a draft
// @Description Discard a draft - status set to rejected, no transaction created
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft UUID"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /drafts/{id}/reject [post]
func (h *DraftHandler) Reject(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "draft")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Draft", "Reject")
	defer span.End()

	if r.Method != http.MethodPost {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		log.Warn("invalid draft id")
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid draft id format", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	log.Info("rejecting draft", "id", id)
	err := h.svc.Reject(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			log.Info("draft not found", "id", id)
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, entities.ErrDraftNotPending) {
			log.Warn("draft not pending", "id", id, "error", err)
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Error("failed to reject draft", "id", id, "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("draft rejected", "id", id)
	h.writeJSON(w, http.StatusOK, models.ErrorResponse{Error: "rejected"})
}

// Delete handles DELETE /drafts/:id
// @Summary Delete a draft
// @Description Hard delete a rejected draft
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft UUID"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /drafts/{id} [delete]
func (h *DraftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log := observability.NewLogger(r.Context(), "http", "draft")
	ctx, span := observability.StartHandlerSpan(log.Context(), "Draft", "Delete")
	defer span.End()

	if r.Method != http.MethodDelete {
		log.Warn("method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		log.Warn("invalid draft id")
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid draft id format", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	log.Info("deleting draft", "id", id)
	err := h.svc.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			log.Info("draft not found", "id", id)
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, entities.ErrDraftNotRejected) {
			log.Warn("draft not rejected", "id", id, "error", err)
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Error("failed to delete draft", "id", id, "error", err)
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("draft deleted", "id", id)
	h.writeJSON(w, http.StatusOK, models.ErrorResponse{Error: "deleted"})
}

func (h *DraftHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *DraftHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}
