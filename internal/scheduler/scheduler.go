package scheduler

import (
	"context"
	"time"
)

type Job interface {
	Name() string
	Schedule() Schedule
	Run(ctx context.Context) error
}

type Schedule interface {
	NextRun(time.Time) time.Time
}

type IntervalSchedule struct {
	Interval time.Duration
}

func (s IntervalSchedule) NextRun(now time.Time) time.Time {
	return now.Add(s.Interval)
}
