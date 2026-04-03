package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// @title Penster API
// @version 1.0
// @description Expense tracker API with accounts and categories management
// @host localhost:8080
// @BasePath /

type ReadyChecker interface {
	Health(ctx context.Context) error
}

type HealthHandler struct {
	readyChecker ReadyChecker
	version      string
}

func NewHealthHandler(readyChecker ReadyChecker) *HealthHandler {
	return &HealthHandler{
		readyChecker: readyChecker,
		version:      "1.0.0",
	}
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// Health handles GET /health
// @Summary Health check
// @Description Returns the health status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resp := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   h.version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Ready handles GET /ready
// @Summary Readiness check
// @Description Returns the readiness status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if h.readyChecker != nil {
		if err := h.readyChecker.Health(r.Context()); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(HealthResponse{
				Status:    "not ready",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Version:   h.version,
			})
			return
		}
	}

	resp := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   h.version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
