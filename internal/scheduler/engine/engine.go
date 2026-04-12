package engine

import (
	"context"
	"sync"
	"time"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/scheduler"
	"github.com/dimasbaguspm/penster/internal/scheduler/jobs"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type Engine struct {
	jobs          []scheduler.Job
	nextRun       map[scheduler.Job]time.Time
	cfg           *config.Config
	rateCurrencyS *service.RateCurrencyService
	ticker        *time.Ticker
	done          chan struct{}
	wg            sync.WaitGroup
	mu            sync.RWMutex
}

func NewEngine(cfg *config.Config, rateCurrencyS *service.RateCurrencyService) *Engine {
	e := &Engine{
		cfg:           cfg,
		rateCurrencyS: rateCurrencyS,
		nextRun:       make(map[scheduler.Job]time.Time),
	}
	e.jobs = append(e.jobs, jobs.NewRateCurrencyJob(cfg, rateCurrencyS))
	return e
}

func (e *Engine) Start(ctx context.Context) {
	log := observability.NewLogger(ctx, "scheduler", "engine")
	e.ticker = time.NewTicker(1 * time.Second)
	e.done = make(chan struct{})

	log.Info("Starting scheduler", "jobs", len(e.jobs))

	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		for {
			select {
			case <-e.done:
				return
			case <-ctx.Done():
				return
			case tick := <-e.ticker.C:
				e.dispatch(tick)
			}
		}
	}()
}

func (e *Engine) Stop() {
	log := observability.NewLogger(context.Background(), "scheduler", "engine")
	if e.ticker != nil {
		e.ticker.Stop()
	}
	close(e.done)
	e.wg.Wait()
	log.Info("Scheduler stopped")
}

func (e *Engine) dispatch(now time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()

	log := observability.NewLogger(context.Background(), "scheduler", "engine")
	for _, job := range e.jobs {
		nextRun := e.nextRun[job]
		if now.After(nextRun) || now.Equal(nextRun) {
			log.Info("Dispatching job", "job", job.Name())
			e.nextRun[job] = now.Add(job.Schedule().NextRun(now).Sub(now))
			go func(j scheduler.Job) {
				if err := j.Run(context.Background()); err != nil {
					log.Error("Job failed", "job", j.Name(), "error", err)
				}
			}(job)
		}
	}
}
