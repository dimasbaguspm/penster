package engine

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/scheduler"
	"github.com/dimasbaguspm/penster/internal/scheduler/jobs"
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
	e.ticker = time.NewTicker(1 * time.Second)
	e.done = make(chan struct{})

	log.Println("Starting scheduler with", len(e.jobs), "jobs")

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
	if e.ticker != nil {
		e.ticker.Stop()
	}
	close(e.done)
	e.wg.Wait()
	log.Println("Scheduler stopped")
}

func (e *Engine) dispatch(now time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, job := range e.jobs {
		nextRun := e.nextRun[job]
		if now.After(nextRun) || now.Equal(nextRun) {
			log.Printf("Dispatching job: %s", job.Name())
			e.nextRun[job] = now.Add(job.Schedule().NextRun(now).Sub(now))
			go func(j scheduler.Job) {
				if err := j.Run(context.Background()); err != nil {
					log.Printf("Job %s failed: %v", j.Name(), err)
				}
			}(job)
		}
	}
}
