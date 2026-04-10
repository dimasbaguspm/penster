package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/interface/dto"
	"github.com/dimasbaguspm/penster/pkg/response"
)

type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// Summary handles GET /reports/summary
// @Summary Get financial report summary
// @Description Get aggregated financial data including totals by type and category breakdown
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /reports/summary [get]
func (h *ReportHandler) Summary(w http.ResponseWriter, r *http.Request) {
	startDate, endDate, err := dto.ParseReportParams(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	summary, err := h.svc.GetSummary(r.Context(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, response.NewResponse(summary))
}

// ByAccount handles GET /reports/by-account
// @Summary Get spending breakdown by account
// @Description Get aggregated spending data grouped by account
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /reports/by-account [get]
func (h *ReportHandler) ByAccount(w http.ResponseWriter, r *http.Request) {
	startDate, endDate, err := dto.ParseReportParams(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	report, err := h.svc.GetByAccount(r.Context(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, response.NewResponse(report))
}

// ByCategory handles GET /reports/by-category
// @Summary Get spending breakdown by category
// @Description Get aggregated spending data grouped by category
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /reports/by-category [get]
func (h *ReportHandler) ByCategory(w http.ResponseWriter, r *http.Request) {
	startDate, endDate, err := dto.ParseReportParams(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	report, err := h.svc.GetByCategory(r.Context(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		// Debug: log the actual error
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("ByCategory error: %v", err))
		return
	}

	h.writeJSON(w, http.StatusOK, response.NewResponse(report))
}

// Trends handles GET /reports/trends
// @Summary Get time-series trend data
// @Description Get time-series data for charting
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /reports/trends [get]
func (h *ReportHandler) Trends(w http.ResponseWriter, r *http.Request) {
	startDate, endDate, err := dto.ParseReportParams(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	report, err := h.svc.GetTrends(r.Context(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, response.NewResponse(report))
}

func (h *ReportHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *ReportHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.NewErrorResponse(message))
}
