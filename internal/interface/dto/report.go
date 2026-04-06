package dto

import (
	"net/http"
	"time"

	"github.com/dimasbaguspm/penster/internal/domain/entities"
)

type ReportParams struct {
	StartDate string
	EndDate   string
}

func ParseReportParams(r *http.Request) (startDate, endDate time.Time, err error) {
	q := r.URL.Query()

	startDateStr := q.Get("start_date")
	endDateStr := q.Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		return time.Time{}, time.Time{}, entities.ErrInvalidDateRange
	}

	startDate, err = time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, entities.ErrInvalidDateRange
	}

	endDate, err = time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, entities.ErrInvalidDateRange
	}

	if startDate.After(endDate) {
		return time.Time{}, time.Time{}, entities.ErrInvalidDateRange
	}

	return startDate, endDate, nil
}
