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
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.Response
// @Router /drafts [get]
func (h *DraftHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Draft", "List")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := dto.ParseDraftListParams(r)

	drafts, total, err := h.svc.List(ctx, params)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	draftList := make([]models.Draft, 0, len(drafts))
	for _, d := range drafts {
		draftList = append(draftList, *d)
	}

	resp := response.NewPaginatedResponse(draftList, 1, params.PageSize, total)
	h.writeJSON(w, http.StatusOK, resp)
}

// Get handles GET /drafts/:id
// @Summary Get draft by ID
// @Description Get a single draft by its UUID
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /drafts/{id} [get]
func (h *DraftHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Draft", "Get")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	draft, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if draft == nil {
		h.writeError(w, http.StatusNotFound, "draft not found")
		return
	}

	resp := response.NewResponse(*draft)
	h.writeJSON(w, http.StatusOK, resp)
}

// Create handles POST /drafts
// @Summary Create a new draft
// @Description Create a new draft with transaction details
// @Tags drafts
// @Accept json
// @Produce json
// @Param request body models.CreateDraftRequest true "Draft creation request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /drafts [post]
func (h *DraftHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Draft", "Create")
	defer span.End()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateCreateDraftRequest(ctx, &req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	draft, err := h.svc.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, entities.ErrAccountNotFound) ||
			errors.Is(err, entities.ErrCategoryNotFound) ||
			errors.Is(err, entities.ErrTransferAccountNotFound) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := response.NewResponse(*draft)
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
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /drafts/{id} [patch]
func (h *DraftHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Draft", "Update")
	defer span.End()

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	var req models.UpdateDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := dto.ValidateUpdateDraftRequest(ctx, &req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	draft, err := h.svc.Update(ctx, id, &req)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, entities.ErrAccountNotFound) ||
			errors.Is(err, entities.ErrCategoryNotFound) ||
			errors.Is(err, entities.ErrTransferAccountNotFound) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Validate same-account transfer after getting current draft
	if req.TransferAccountID != nil && req.AccountID != nil &&
		*req.TransferAccountID == *req.AccountID {
		h.writeError(w, http.StatusBadRequest, "transfer account cannot be the same as source account")
		return
	}

	resp := response.NewResponse(*draft)
	h.writeJSON(w, http.StatusOK, resp)
}

// Confirm handles POST /drafts/:id/confirm
// @Summary Confirm a draft
// @Description Promote a draft to a committed transaction
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /drafts/{id}/confirm [post]
func (h *DraftHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Draft", "Confirm")
	defer span.End()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	tx, err := h.svc.Confirm(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, entities.ErrDraftNotPending) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, entities.ErrInsufficientBalance) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := response.NewResponse(*tx)
	h.writeJSON(w, http.StatusOK, resp)
}

// Reject handles POST /drafts/:id/reject
// @Summary Reject a draft
// @Description Discard a draft - status set to rejected, no transaction created
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /drafts/{id}/reject [post]
func (h *DraftHandler) Reject(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Draft", "Reject")
	defer span.End()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	err := h.svc.Reject(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, entities.ErrDraftNotPending) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, response.NewResponse(map[string]string{"status": "rejected"}))
}

// Delete handles DELETE /drafts/:id
// @Summary Delete a draft
// @Description Hard delete a rejected draft
// @Tags drafts
// @Accept json
// @Produce json
// @Param id path string true "Draft UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /drafts/{id} [delete]
func (h *DraftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.StartHandlerSpan(r.Context(), "Draft", "Delete")
	defer span.End()

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid draft id")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid draft id format")
		return
	}

	err := h.svc.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrDraftNotFound) {
			h.writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, entities.ErrDraftNotRejected) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, response.NewResponse(map[string]string{"status": "deleted"}))
}

func (h *DraftHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *DraftHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.NewErrorResponse(message))
}
