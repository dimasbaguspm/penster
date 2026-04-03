package engine

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/service"
)

type Engine struct {
	jobs     []Job
	nextRun  map[Job]time.Time
	cfg      *config.Config
	svc      *service.RateCurrencyService
	ticker   *time.Ticker
	done     chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

func NewEngine(cfg *config.Config, svc *service.RateCurrencyService) *Engine {
	return &Engine{
		cfg:     cfg,
		svc:     svc,
		nextRun: make(map[Job]time.Time),
	}
}

func (e *Engine) RegisterJob(job Job) {
	e.jobs = append(e.jobs, job)
	e.nextRun[job] = time.Now()
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
			go func(j Job) {
				if err := j.Run(context.Background(), e.svc); err != nil {
					log.Printf("Job %s failed: %v", j.Name(), err)
				}
			}(job)
		}
	}
}
