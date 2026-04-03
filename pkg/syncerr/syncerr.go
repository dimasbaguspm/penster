package syncerr

import (
	"sync"
)

// Group manages goroutines and collects errors.
type Group struct {
	wg   sync.WaitGroup
	errs []error
	mu   sync.Mutex
}

// Go runs fn in a goroutine and collects any error returned.
func (g *Group) Go(fn func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := fn(); err != nil {
			g.mu.Lock()
			g.errs = append(g.errs, err)
			g.mu.Unlock()
		}
	}()
}

// Wait blocks until all goroutines complete and returns collected errors.
func (g *Group) Wait() []error {
	g.wg.Wait()
	return g.errs
}
