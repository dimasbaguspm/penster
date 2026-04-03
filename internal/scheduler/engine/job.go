package engine

import (
	"context"
	"time"

	"github.com/dimasbaguspm/penster/internal/application/service"
)

type Job interface {
	Name() string
	Schedule() Schedule
	Run(ctx context.Context, svc *service.RateCurrencyService) error
}

type Schedule interface {
	NextRun(now time.Time) time.Time
}

type IntervalSchedule struct {
	Interval time.Duration
}

func (s IntervalSchedule) NextRun(now time.Time) time.Time {
	return now.Add(s.Interval)
}
